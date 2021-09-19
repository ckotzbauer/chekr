package resources

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

func (r Resource) Execute() printer.PrintableList {
	pods := r.KubeClient.ListPods(kubernetes.PodQuery{
		Namespace:          r.Namespace,
		LabelSelector:      r.LabelSelector,
		Names:              r.Pods,
		AnnotationSelector: r.AnnotationSelector,
	})

	if len(pods) == 0 {
		return PodValuesList{}
	}

	if kubernetes.IsResourceName(r.Prometheus.Url) {
		readyChannel := make(chan struct{})
		stopChannel := make(chan struct{}, 1)

		fn := func() PodValuesList {
			<-readyChannel
			x := r.executeInternal(pods)
			close(stopChannel)
			return x
		}

		ch := channelify.Channelify(fn)
		c := ch.(func() chan PodValuesList)()
		r.KubeClient.ForwardResource(&r.Prometheus, readyChannel, stopChannel)

		result := <-c
		return result
	} else {
		return r.executeInternal(pods)
	}
}

func (r Resource) executeInternal(pods []corev1.Pod) PodValuesList {
	podValuesList := PodValuesList{Items: []PodValues{}}

	v1api := r.Prometheus.InitPrometheus()
	queryRange := v1.Range{
		Start: time.Now().Add(-time.Hour * 24 * time.Duration(r.Prometheus.CountDays)),
		End:   time.Now(),
		Step:  time.Minute * 5,
	}

	fn1 := func(r Resource, pod corev1.Pod, v1api v1.API, queryRange v1.Range) printer.Printable {
		return r.analyzePod(pod, v1api, queryRange)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.Printable)

	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning {
			ch := ch1.(func(Resource, corev1.Pod, v1.API, v1.Range) chan printer.Printable)(r, pod, v1api, queryRange)
			channels = append(channels, ch)
		}
	}

	for _, v := range channels {
		result := <-v

		podValuesList.Items = append(podValuesList.Items, result.(PodValues))
	}

	return podValuesList
}

func (r Resource) analyzePod(pod corev1.Pod, v1api v1.API, queryRange v1.Range) printer.Printable {
	podValues := PodValues{
		Namespace:  pod.Namespace,
		Pod:        pod.Name,
		Containers: []ContainerValue{},
	}

	matrix, err := queryMatrix(r.Prometheus, v1api, MetricsQuery(pod.Namespace, pod.Name), queryRange)

	if err != nil {
		logrus.WithError(err).WithField("pod", pod.Namespace+"/"+pod.Name).Fatalf("Could not query metrics for pod!")
	}

	for _, container := range pod.Spec.Containers {
		cv := ContainerValue{
			Name:           container.Name,
			MemoryRequests: AnalyzedValues{},
			MemoryLimits:   AnalyzedValues{},
			CPURequests:    AnalyzedValues{},
			CPULimits:      AnalyzedValues{},
		}

		calculate(
			MemoryUsageMetric(matrix, container.Name),
			container.Resources.Requests.Memory(),
			container.Resources.Limits.Memory(),
			&cv.MemoryRequests,
			&cv.MemoryLimits)

		calculate(
			CPUUsageMetric(matrix, container.Name),
			container.Resources.Requests.Cpu(),
			container.Resources.Limits.Cpu(),
			&cv.CPURequests,
			&cv.CPULimits)

		podValues.Containers = append(podValues.Containers, cv)
	}

	return podValues
}

func calculate(usageMetric *model.SampleStream, requests, limits *resource.Quantity, analyzedRequests, analyzedLimits *AnalyzedValues) {
	var usages []float64

	if usageMetric == nil {
		return
	}

	for _, usg := range usageMetric.Values {
		usages = append(usages, float64(usg.Value))
	}

	analyzedRequests.Min = util.ComputedValue{Value: util.MinOf(usages...)}
	analyzedRequests.Max = util.ComputedValue{Value: util.MaxOf(usages...)}
	analyzedRequests.Avg = util.ComputedValue{Value: util.SumOf(usages...) / float64(len(usageMetric.Values))}

	if !requests.IsZero() {
		analyzedRequests.Min.Percentage = analyzedRequests.Min.Value / requests.AsApproximateFloat64()
		analyzedRequests.Max.Percentage = analyzedRequests.Max.Value / requests.AsApproximateFloat64()
		analyzedRequests.Avg.Percentage = analyzedRequests.Avg.Value / requests.AsApproximateFloat64()
		analyzedRequests.Current = requests.AsApproximateFloat64()
		analyzedRequests.HasValue = true
	}

	analyzedLimits.Min = util.ComputedValue{Value: analyzedRequests.Min.Value}
	analyzedLimits.Max = util.ComputedValue{Value: analyzedRequests.Max.Value}
	analyzedLimits.Avg = util.ComputedValue{Value: analyzedRequests.Avg.Value}

	if !limits.IsZero() {
		analyzedLimits.Min.Percentage = analyzedLimits.Min.Value / limits.AsApproximateFloat64()
		analyzedLimits.Max.Percentage = analyzedLimits.Max.Value / limits.AsApproximateFloat64()
		analyzedLimits.Avg.Percentage = analyzedLimits.Avg.Value / limits.AsApproximateFloat64()
		analyzedLimits.Current = limits.AsApproximateFloat64()
		analyzedLimits.HasValue = true
	}
}

func queryMatrix(prom prometheus.Prometheus, v1api v1.API, query string, r v1.Range) (model.Matrix, error) {
	value, err := prom.QueryRange(v1api, query, r)

	if err != nil {
		return nil, err
	}

	matrix := value.(model.Matrix)
	return matrix, err
}
