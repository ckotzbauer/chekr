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
	LimitsThreshold    int
	RequestsThreshold  int
}

type AnalyzedValues struct {
	Min      util.ComputedValue
	Max      util.ComputedValue
	Avg      util.ComputedValue
	Current  util.ComputedValue
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
					c.MemoryRequests.Current.FormatMemory(),
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
					c.MemoryLimits.Current.FormatMemory(),
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
					c.CPURequests.Current.FormatCPU(),
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
					c.CPULimits.Current.FormatCPU(),
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
