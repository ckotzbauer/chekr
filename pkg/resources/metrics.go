package resources

import (
	"fmt"

	"github.com/prometheus/common/model"
)

const (
	memoryUsageMetric = "container_memory_working_set_bytes"
	cpuUsageMetric    = "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate"
)

type convert func(*model.SampleStream) bool

func findMetric(matrix model.Matrix, name string, fn convert) *model.SampleStream {
	fullStream := model.SampleStream{Values: []model.SamplePair{}}

	for _, v := range matrix {
		if string(v.Metric["__name__"]) == name && fn(v) {
			fullStream.Values = append(fullStream.Values, v.Values...)
			fullStream.Metric = v.Metric
		}
	}

	if len(fullStream.Values) > 0 {
		return &fullStream
	}

	return nil
}

func MemoryUsageMetric(matrix model.Matrix, container string) *model.SampleStream {
	return findMetric(matrix, memoryUsageMetric, func(ss *model.SampleStream) bool {
		return ss.Metric["container"] == model.LabelValue(container)
	})
}

func CPUUsageMetric(matrix model.Matrix, container string) *model.SampleStream {
	return findMetric(matrix, cpuUsageMetric, func(ss *model.SampleStream) bool {
		return ss.Metric["container"] == model.LabelValue(container)
	})
}

func MetricsQuery(namespace, pod string) string {
	return fmt.Sprintf("{__name__=~\"node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate|container_memory_working_set_bytes\", namespace=\"%v\", pod=\"%v\", container!=\"\", container!=\"POD\"}", namespace, pod)
}
