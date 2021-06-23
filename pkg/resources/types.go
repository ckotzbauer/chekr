package resources

import (
	"bytes"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
)

type Resource struct {
	Prometheus         prometheus.Prometheus
	KubeOverrides      *clientcmd.ConfigOverrides
	KubeClient         *kubernetes.KubeClient
	Pods               []string
	Namespace          string
	LabelSelector      string
	AnnotationSelector string
}

type AnalyzedValues struct {
	Min util.ComputedValue
	Max util.ComputedValue
	Avg util.ComputedValue
}

type PodValues struct {
	Namespace      string
	Pod            string
	MemoryRequests AnalyzedValues
	MemoryLimits   AnalyzedValues
	CPURequests    AnalyzedValues
	CPULimits      AnalyzedValues
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
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{
		"Namespace",
		"Pod",
		"Memory Requests",
		"Memory Limits",
		"CPU Requests",
		"CPU Limits",
	})

	for _, v := range p.Items {
		table.Append([]string{
			v.Namespace,
			v.Pod,
			v.MemoryRequests.Min.FormatMemory() + "\n" + v.MemoryRequests.Avg.FormatMemory() + "\n" + v.MemoryRequests.Max.FormatMemory(),
			v.MemoryLimits.Min.FormatMemory() + "\n" + v.MemoryLimits.Avg.FormatMemory() + "\n" + v.MemoryLimits.Max.FormatMemory(),
			v.CPURequests.Min.FormatCPU() + "\n" + v.CPURequests.Avg.FormatCPU() + "\n" + v.CPURequests.Max.FormatCPU(),
			v.CPULimits.Min.FormatCPU() + "\n" + v.CPULimits.Avg.FormatCPU() + "\n" + v.CPULimits.Max.FormatCPU(),
		})
	}

	table.Render()
	return buf.String()
}
