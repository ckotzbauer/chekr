package resources

import (
	"fmt"

	"github.com/prometheus/common/model"
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

func (r Resource) MemoryUsageMetric(matrix model.Matrix, container string) *model.SampleStream {
	return findMetric(matrix, r.MemoryMetric, func(ss *model.SampleStream) bool {
		return ss.Metric["container"] == model.LabelValue(container)
	})
}

func (r Resource) CPUUsageMetric(matrix model.Matrix, container string) *model.SampleStream {
	return findMetric(matrix, r.CpuMetric, func(ss *model.SampleStream) bool {
		return ss.Metric["container"] == model.LabelValue(container)
	})
}

func (r Resource) MetricsQuery(namespace, pod string) string {
	return fmt.Sprintf("{__name__=~\"%v|%v\", namespace=\"%v\", pod=\"%v\", container!=\"\", container!=\"POD\"}", r.CpuMetric, r.MemoryMetric, namespace, pod)
}
