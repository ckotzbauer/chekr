package ha

import (
	"strings"

	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (h HighAvailability) Execute() (printer.PrintableList, error) {
	var pods []corev1.Pod

	if h.Selector == "" {
		pods = h.KubeClient.GetNamespacedPods(h.Namespace, h.Pods)
	} else if h.Selector != "" {
		pods = h.KubeClient.ListPods(h.Namespace, h.Selector)
	}

	fn1 := func(h HighAvailability, pod corev1.Pod) printer.PrintableResult {
		return h.analyzePod(pod)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.PrintableResult)
	podAvailabilityList := PodAvailabilityList{Items: []PodAvailability{}}

	for _, pod := range pods {
		ch := ch1.(func(HighAvailability, corev1.Pod) chan printer.PrintableResult)(h, pod)
		channels = append(channels, ch)
	}

	for _, v := range channels {
		result := <-v

		if result.Error != nil {
			return nil, result.Error
		}

		if result.Item == nil {
			continue
		}

		podAvailabilityList.Items = append(podAvailabilityList.Items, result.Item.(PodAvailability))
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

	return uniqueList, nil
}

func (h HighAvailability) analyzePod(pod corev1.Pod) printer.PrintableResult {
	availability := PodAvailability{
		Namespace: pod.Namespace,
	}

	kind := util.GetOwnerKind(pod.OwnerReferences)
	var err error

	if kind == "ReplicaSet" {
		err = h.analyzeReplicaSet(&availability, pod)
	} else if kind == "StatefulSet" {
		err = h.analyzeStatefulSet(pod.OwnerReferences, &availability)
	} else if kind == "DaemonSet" {
		err = h.analyzeDaemonSet(pod.OwnerReferences, &availability)
	} else if kind == "Job" {
		return printer.PrintableResult{}
	} else if kind != "" {
		// A CRD-Operator created the pod
		availability.Type = kind
		availability.Replicas = 1
	} else {
		// Single pod without owner
		availability.Type = "Pod"
		availability.Replicas = 1
	}

	if err != nil {
		return printer.PrintableResult{Error: err}
	}

	return printer.PrintableResult{Item: availability}
}

func (h HighAvailability) analyzeReplicaSet(availability *PodAvailability, pod corev1.Pod) error {
	rs, err := h.KubeClient.GetReplicaSet(pod)

	if err != nil {
		return err
	}

	availability.Type = util.GetOwnerKind(rs.OwnerReferences)
	availability.Replicas = *rs.Spec.Replicas

	if availability.Type == "Deployment" {
		err = h.analyzeDeployment(rs.OwnerReferences, availability)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h HighAvailability) analyzeDeployment(refs []metav1.OwnerReference, availability *PodAvailability) error {
	deployment, err := h.KubeClient.GetDeployment(availability.Namespace, refs[0].Name)

	if err != nil {
		return err
	}

	availability.Name = deployment.Name
	availability.Owner = "deployment/" + deployment.Name
	availability.RolloutStrategy = string(deployment.Spec.Strategy.Type)
	return h.analyzePodTemplace(deployment.Spec.Template, availability)
}

func (h HighAvailability) analyzeDaemonSet(refs []metav1.OwnerReference, availability *PodAvailability) error {
	daemonSet, err := h.KubeClient.GetDaemonSet(availability.Namespace, refs[0].Name)

	if err != nil {
		return err
	}

	availability.Name = daemonSet.Name
	availability.Type = "DaemonSet"
	availability.Replicas = daemonSet.Status.NumberReady
	availability.Owner = "daemonset/" + daemonSet.Name
	availability.RolloutStrategy = string(daemonSet.Spec.UpdateStrategy.Type)
	return h.analyzePodTemplace(daemonSet.Spec.Template, availability)
}

func (h HighAvailability) analyzeStatefulSet(refs []metav1.OwnerReference, availability *PodAvailability) error {
	statefulSet, err := h.KubeClient.GetStatefulSet(availability.Namespace, refs[0].Name)

	if err != nil {
		return err
	}

	availability.Type = "StatefulSet"
	availability.Owner = "statefulset/" + statefulSet.Name
	availability.Replicas = *statefulSet.Spec.Replicas
	availability.Name = statefulSet.Name
	availability.RolloutStrategy = string(statefulSet.Spec.UpdateStrategy.Type)

	err = h.analyzePodTemplace(statefulSet.Spec.Template, availability)

	if err != nil {
		return err
	}

	if len(statefulSet.Spec.VolumeClaimTemplates) > 0 {
		if availability.PVC == "" {
			availability.PVC = "STS"
		} else {
			availability.PVC += ",STS"
		}
	}

	return nil
}

func (h HighAvailability) analyzePodTemplace(spec corev1.PodTemplateSpec, availability *PodAvailability) error {
	if spec.Spec.Affinity != nil && spec.Spec.Affinity.PodAntiAffinity != nil {
		affinity := spec.Spec.Affinity.PodAntiAffinity
		matched, err := hasMatchingPodAntiAffinity(affinity, spec.Labels)

		if err != nil {
			return err
		}

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
				return err
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
	return nil
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

func hasMatchingPodAntiAffinity(affinity *corev1.PodAntiAffinity, podLabels map[string]string) (bool, error) {
	for _, reqAffinity := range affinity.RequiredDuringSchedulingIgnoredDuringExecution {
		selector, err := metav1.LabelSelectorAsSelector(reqAffinity.LabelSelector)

		if err != nil {
			return false, err
		}

		if selector.Matches(labels.Set(podLabels)) {
			return true, nil
		}
	}

	return false, nil
}
