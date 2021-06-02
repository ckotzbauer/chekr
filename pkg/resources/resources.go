package resources

import (
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func (r Resource) Execute() (printer.PrintableList, error) {
	v1api := r.Prometheus.InitPrometheus()
	queryRange := v1.Range{
		Start: time.Now().Add(-time.Hour * 24 * 30),
		End:   time.Now(),
		Step:  time.Minute * 5,
	}

	var pods []corev1.Pod

	if r.Selector == "" {
		pods = r.KubeClient.GetNamespacedPods(r.Namespace, r.Pods)
	} else if r.Selector != "" {
		pods = r.KubeClient.ListPods(r.Namespace, r.Selector)
	}

	fn1 := func(r Resource, pod corev1.Pod, v1api v1.API, queryRange v1.Range) printer.PrintableResult {
		return r.analyzePod(pod, v1api, queryRange)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.PrintableResult)
	podValuesList := PodValuesList{Items: []PodValues{}}

	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning {
			ch := ch1.(func(Resource, corev1.Pod, v1.API, v1.Range) chan printer.PrintableResult)(r, pod, v1api, queryRange)
			channels = append(channels, ch)
		}
	}

	for _, v := range channels {
		result := <-v

		if result.Error != nil {
			return nil, result.Error
		}

		podValuesList.Items = append(podValuesList.Items, result.Item.(PodValues))
	}

	return podValuesList, nil
}

func (r Resource) analyzePod(pod corev1.Pod, v1api v1.API, queryRange v1.Range) printer.PrintableResult {
	podValues := PodValues{
		Namespace:      pod.Namespace,
		Pod:            pod.Name,
		MemoryRequests: AnalyzedValues{},
		MemoryLimits:   AnalyzedValues{},
		CPURequests:    AnalyzedValues{},
		CPULimits:      AnalyzedValues{},
	}

	matrix, err := queryMatrix(r.Prometheus, v1api, Memory(pod.Namespace, pod.Name), queryRange)

	if err != nil {
		return printer.PrintableResult{Error: err}
	}

	calculate(
		findMetric(matrix, MemoryUsageMetric),
		findMetric(matrix, MemoryRequestsMetric),
		findMetric(matrix, MemoryLimitsMetric),
		&podValues.MemoryRequests,
		&podValues.MemoryLimits)

	calculate(
		findMetric(matrix, CPUUsageMetric),
		findMetric(matrix, CPURequestsMetric),
		findMetric(matrix, CPULimitsMetric),
		&podValues.CPURequests,
		&podValues.CPULimits)

	return printer.PrintableResult{Item: podValues}
}

func calculate(usageMetric, requestMetric, limitMetric *model.SampleStream, analyzedRequests, analyzedLimits *AnalyzedValues) {
	var avgRequests float64
	var avgRequestCounter float64
	var avgLimits float64
	var avgLimitCounter float64

	for i := 0; i < len(usageMetric.Values); i++ {
		usg := usageMetric.Values[i]
		req := findPair(requestMetric.Values, usg.Timestamp)
		lim := findPair(limitMetric.Values, usg.Timestamp)

		var requestPercentage float64
		var limitPercentage float64

		if req.Value != 0 {
			requestPercentage = float64(usg.Value) / float64(req.Value)
			avgRequests += requestPercentage
			avgRequestCounter++
		}

		if lim.Value != 0 {
			limitPercentage = float64(usg.Value) / float64(lim.Value)
			avgLimits += limitPercentage
			avgLimitCounter++
		}

		analyzedRequests.Min = computeMin(analyzedRequests.Min, requestPercentage, float64(usg.Value))
		analyzedLimits.Min = computeMin(analyzedLimits.Min, limitPercentage, float64(usg.Value))
		analyzedRequests.Max = computeMax(analyzedRequests.Max, requestPercentage, float64(usg.Value))
		analyzedLimits.Max = computeMax(analyzedLimits.Max, limitPercentage, float64(usg.Value))
	}

	if avgRequestCounter > 0 {
		avg := avgRequests / avgRequestCounter
		// TODO: fix average value
		analyzedRequests.Avg = util.ComputedValue{Percentage: avg, Value: 0}
	}

	if avgLimitCounter > 0 {
		avg := avgLimits / avgLimitCounter
		// TODO: fix average value
		analyzedLimits.Avg = util.ComputedValue{Percentage: avg, Value: 0}
	}
}

func findMetric(matrix model.Matrix, metric string) *model.SampleStream {
	var potential []*model.SampleStream

	for i := 0; i < matrix.Len(); i++ {
		if string(matrix[i].Metric["__name__"]) == metric {
			potential = append(potential, matrix[i])
		}
	}

	for i := 0; i < len(potential); i++ {
		if len(potential[i].Values) > 0 && float64(potential[i].Values[len(potential[i].Values)-1].Value) > 0 {
			return potential[i]
		}
	}

	return nil
}

func findPair(pairs []model.SamplePair, time model.Time) model.SamplePair {
	for i := 0; i < len(pairs); i++ {
		p := pairs[i]
		if p.Timestamp == time {
			return p
		}
	}

	return model.SamplePair{}
}

func queryMatrix(prom prometheus.Prometheus, v1api v1.API, query string, r v1.Range) (model.Matrix, error) {
	value, err := prom.QueryRange(v1api, query, r)

	if err != nil {
		return nil, err
	}

	matrix := value.(model.Matrix)
	return matrix, err
}

func computeMin(computed util.ComputedValue, current, value float64) util.ComputedValue {
	if computed.Percentage < current && computed.Percentage != 0 {
		return computed
	}

	return util.ComputedValue{Value: value, Percentage: current}
}

func computeMax(computed util.ComputedValue, current, value float64) util.ComputedValue {
	if computed.Percentage > current {
		return computed
	}

	return util.ComputedValue{Value: value, Percentage: current}
}
