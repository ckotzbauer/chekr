package ha

import (
	"strings"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (h HighAvailability) Execute() printer.PrintableList {
	pods := h.KubeClient.ListPods(kubernetes.PodQuery{
		Namespace:          h.Namespace,
		LabelSelector:      h.LabelSelector,
		Names:              h.Pods,
		AnnotationSelector: h.AnnotationSelector,
	})

	fn1 := func(h HighAvailability, pod corev1.Pod) printer.Printable {
		return h.analyzePod(pod)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.Printable)
	podAvailabilityList := PodAvailabilityList{Items: []PodAvailability{}}

	for _, pod := range pods {
		ch := ch1.(func(HighAvailability, corev1.Pod) chan printer.Printable)(h, pod)
		channels = append(channels, ch)
	}

	for _, v := range channels {
		result := <-v

		if result == nil {
			continue
		}

		podAvailabilityList.Items = append(podAvailabilityList.Items, result.(PodAvailability))
	}

	owners := make([]string, 0)
	uniqueList := PodAvailabilityList{Items: make([]PodAvailability, 0)}

	for _, i := range podAvailabilityList.Items {
		if util.Contains(owners, i.Owner) {
			continue
		}

		owners = append(owners, i.Owner)
		h.rankPod(&i)
		uniqueList.Items = append(uniqueList.Items, i)
	}

	return uniqueList
}

func (h HighAvailability) analyzePod(pod corev1.Pod) printer.Printable {
	availability := PodAvailability{
		Namespace: pod.Namespace,
	}

	kind := util.GetOwnerKind(pod.OwnerReferences)

	if kind == "ReplicaSet" {
		h.analyzeReplicaSet(&availability, pod)
	} else if kind == "StatefulSet" {
		h.analyzeStatefulSet(pod.OwnerReferences, &availability)
	} else if kind == "DaemonSet" {
		h.analyzeDaemonSet(pod.OwnerReferences, &availability)
	} else if kind == "Job" {
		return availability
	} else if kind != "" {
		// A CRD-Operator created the pod
		availability.Type = kind
		availability.Replicas = 1
	} else {
		// Single pod without owner
		availability.Type = "Pod"
		availability.Replicas = 1
	}

	return availability
}

func (h HighAvailability) analyzeReplicaSet(availability *PodAvailability, pod corev1.Pod) {
	rs, err := h.KubeClient.GetReplicaSet(pod)

	if err != nil {
		logrus.WithError(err).WithField("pod", pod.Namespace+"/"+pod.Name).Fatalf("Could not get ReplicaSet for pod!")
	}

	availability.Type = util.GetOwnerKind(rs.OwnerReferences)
	availability.Replicas = *rs.Spec.Replicas

	if availability.Type == "Deployment" {
		h.analyzeDeployment(rs.OwnerReferences, availability)
	}
}

func (h HighAvailability) analyzeDeployment(refs []metav1.OwnerReference, availability *PodAvailability) {
	deployment, err := h.KubeClient.GetDeployment(availability.Namespace, refs[0].Name)

	if err != nil {
		logrus.WithError(err).WithField("deployment", availability.Namespace+"/"+refs[0].Name).Fatalf("Could not get Deployment!")
	}

	availability.Name = deployment.Name
	availability.Owner = "deployment/" + deployment.Name
	availability.RolloutStrategy = string(deployment.Spec.Strategy.Type)
	h.analyzePodTemplace(deployment.Spec.Template, availability)
}

func (h HighAvailability) analyzeDaemonSet(refs []metav1.OwnerReference, availability *PodAvailability) {
	daemonSet, err := h.KubeClient.GetDaemonSet(availability.Namespace, refs[0].Name)

	if err != nil {
		logrus.WithError(err).WithField("daemonset", availability.Namespace+"/"+refs[0].Name).Fatalf("Could not get Daemonset for pod!")
	}

	availability.Name = daemonSet.Name
	availability.Type = "DaemonSet"
	availability.Replicas = daemonSet.Status.NumberReady
	availability.Owner = "daemonset/" + daemonSet.Name
	availability.RolloutStrategy = string(daemonSet.Spec.UpdateStrategy.Type)
	h.analyzePodTemplace(daemonSet.Spec.Template, availability)
}

func (h HighAvailability) analyzeStatefulSet(refs []metav1.OwnerReference, availability *PodAvailability) {
	statefulSet, err := h.KubeClient.GetStatefulSet(availability.Namespace, refs[0].Name)

	if err != nil {
		logrus.WithError(err).WithField("statefulset", availability.Namespace+"/"+refs[0].Name).Fatalf("Could not get StatefulSet for pod!")
	}

	availability.Type = "StatefulSet"
	availability.Owner = "statefulset/" + statefulSet.Name
	availability.Replicas = *statefulSet.Spec.Replicas
	availability.Name = statefulSet.Name
	availability.RolloutStrategy = string(statefulSet.Spec.UpdateStrategy.Type)

	h.analyzePodTemplace(statefulSet.Spec.Template, availability)

	if len(statefulSet.Spec.VolumeClaimTemplates) > 0 {
		if availability.PVC == "" {
			availability.PVC = "STS"
		} else {
			availability.PVC += ",STS"
		}
	}
}

func (h HighAvailability) analyzePodTemplace(spec corev1.PodTemplateSpec, availability *PodAvailability) {
	if spec.Spec.Affinity != nil && spec.Spec.Affinity.PodAntiAffinity != nil {
		affinity := spec.Spec.Affinity.PodAntiAffinity
		matched := hasMatchingPodAntiAffinity(affinity, spec.Labels)

		if matched {
			availability.PodAntiAffinity = "Yes"
		} else {
			availability.PodAntiAffinity = "No"
		}
	} else {
		availability.PodAntiAffinity = "No"
	}

	pvcs := []string{}

	for _, volume := range spec.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName != "" {
			pvc, err := h.KubeClient.GetPersistentVolumeClaim(h.Namespace, volume.PersistentVolumeClaim.ClaimName)

			if err != nil {
				logrus.WithError(err).WithField("pvc", h.Namespace+"/"+volume.PersistentVolumeClaim.ClaimName).Fatalf("Could not get PVC!")
			}

			for _, am := range pvc.Spec.AccessModes {
				if am == corev1.ReadOnlyMany {
					pvcs = append(pvcs, "ROX")
				} else if am == corev1.ReadWriteOnce {
					pvcs = append(pvcs, "RWO")
				} else if am == corev1.ReadWriteMany {
					pvcs = append(pvcs, "RWX")
				}
			}
		}
	}

	availability.PVC = strings.Join(pvcs, ",")
}

func (h HighAvailability) rankPod(pod *PodAvailability) {
	// Ranks
	// 0 undefinied
	// 1 high-available (failure-resilient, zero-downtime-deployment capable)
	// 2 zero-downtime-deployment capable (non failure-resilient)
	// 3 single-point-of-failure
	// 4 standalone pod

	if pod.Type == "Pod" {
		pod.Rank = 4
	}

	if pod.Replicas == 1 {
		if pod.RolloutStrategy != "RollingUpdate" || pod.PVC == "RWO" {
			pod.Rank = 3
		} else {
			pod.Rank = 2
		}
	} else {
		if pod.PodAntiAffinity == "No" {
			pod.Rank = 2
		}

		if pod.PodAntiAffinity == "Yes" {
			pod.Rank = 1
		}

		if pod.RolloutStrategy == "Recreate" {
			pod.Rank = 3
		}
	}
}

func hasMatchingPodAntiAffinity(affinity *corev1.PodAntiAffinity, podLabels map[string]string) bool {
	for _, reqAffinity := range affinity.RequiredDuringSchedulingIgnoredDuringExecution {
		selector, err := metav1.LabelSelectorAsSelector(reqAffinity.LabelSelector)

		if err != nil {
			logrus.WithError(err).WithField("selector", reqAffinity.LabelSelector).Fatalf("Could not convert label-selector!")
		}

		if selector.Matches(labels.Set(podLabels)) {
			return true
		}
	}

	return false
}
