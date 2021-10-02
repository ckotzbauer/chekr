package resources

import (
	"bytes"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/olekukonko/tablewriter"
)

type Resource struct {
	Prometheus         prometheus.Prometheus
	KubeClient         *kubernetes.KubeClient
	Pods               []string
	Namespace          string
	LabelSelector      string
	AnnotationSelector string
	CpuMetric          string
	MemoryMetric       string
}

type AnalyzedValues struct {
	Min      util.ComputedValue
	Max      util.ComputedValue
	Avg      util.ComputedValue
	Current  float64
	HasValue bool
}

type ContainerValue struct {
	Name           string
	MemoryRequests AnalyzedValues
	MemoryLimits   AnalyzedValues
	CPURequests    AnalyzedValues
	CPULimits      AnalyzedValues
}

type PodValues struct {
	Namespace  string
	Pod        string
	Containers []ContainerValue
}

// PodValuesList implements PrintableList
type PodValuesList struct {
	Items []PodValues
}

func (p PodValuesList) ToJson() string {
	return printer.ToJson(p.Items)
}

func (p PodValuesList) ToHtml() string {
	return printer.ToHtml(HtmlPage, p)
}

func (p PodValuesList) ToTable() string {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{
		"Namespace",
		"Pod",
		"Container",
		"",
		"Current value",
		"Min value",
		"Average value",
		"Max value",
	})

	for _, v := range p.Items {
		for _, c := range v.Containers {
			if c.MemoryRequests.HasValue {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"Memory Requests",
					util.ByteCountIEC(c.MemoryRequests.Current),
					c.MemoryRequests.Min.FormatMemory(),
					c.MemoryRequests.Avg.FormatMemory(),
					c.MemoryRequests.Max.FormatMemory(),
				})
			} else {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"Memory",
					"",
					c.MemoryRequests.Min.FormatMemory(),
					c.MemoryRequests.Avg.FormatMemory(),
					c.MemoryRequests.Max.FormatMemory(),
				})
			}

			if c.MemoryLimits.HasValue {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"Memory Limits",
					util.ByteCountIEC(c.MemoryLimits.Current),
					c.MemoryLimits.Min.FormatMemory(),
					c.MemoryLimits.Avg.FormatMemory(),
					c.MemoryLimits.Max.FormatMemory(),
				})
			} else {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"Memory",
					"",
					c.MemoryLimits.Min.FormatMemory(),
					c.MemoryLimits.Avg.FormatMemory(),
					c.MemoryLimits.Max.FormatMemory(),
				})
			}

			if c.CPURequests.HasValue {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"CPU Requests",
					util.Cores(c.CPURequests.Current),
					c.CPURequests.Min.FormatCPU(),
					c.CPURequests.Avg.FormatCPU(),
					c.CPURequests.Max.FormatCPU(),
				})
			} else {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"CPUs",
					"",
					c.CPURequests.Min.FormatCPU(),
					c.CPURequests.Avg.FormatCPU(),
					c.CPURequests.Max.FormatCPU(),
				})
			}

			if c.CPULimits.HasValue {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"CPU Limits",
					util.Cores(c.CPULimits.Current),
					c.CPULimits.Min.FormatCPU(),
					c.CPULimits.Avg.FormatCPU(),
					c.CPULimits.Max.FormatCPU(),
				})
			} else {
				table.Append([]string{
					v.Namespace,
					v.Pod,
					c.Name,
					"CPUs",
					"",
					c.CPULimits.Min.FormatCPU(),
					c.CPULimits.Avg.FormatCPU(),
					c.CPULimits.Max.FormatCPU(),
				})
			}
		}

	}

	table.Render()
	return buf.String()
}
