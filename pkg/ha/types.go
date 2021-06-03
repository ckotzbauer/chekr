package ha

import (
	"bytes"
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
)

var description = `

Ranks
0: Undefined (detection was not possible)
1: High-available (failure-resilient, zero-downtime-deployment capable)
2: Zero-downtime-deployment capable: (non failure-resilient)
3: Single-point-of-failure
4: Standalone pod

Definitions
Failure resilient: If a single pod or node crashes, the availability of the workload is not affected.
Zero-downtime-deployment: The workload can be updated without downtime
Single-point-of-failure: If a single pod or the node crashes, the availability of the workload is degraded.
Standalone pod: If the pod is removed or crashes, it disappears without replacement.
`

type HighAvailability struct {
	KubeOverrides *clientcmd.ConfigOverrides
	KubeClient    *kubernetes.KubeClient
	Pods          []string
	Namespace     string
	Selector      string
}

type PodAvailability struct {
	Namespace       string
	Name            string
	Owner           string
	Type            string
	Replicas        int32
	PodAntiAffinity string
	RolloutStrategy string
	PVC             string
	Rank            int32
}

// PodAvailabilityList implements PrintableList
type PodAvailabilityList struct {
	Items []PodAvailability
}

func (p PodAvailabilityList) ToJson() string {
	return printer.ToJson(p.Items)
}

func (p PodAvailabilityList) ToHtml() string {
	return printer.ToHtml(fmt.Sprintf(HtmlPage, description), p)
}

func (p PodAvailabilityList) ToTable() string {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{
		"Namespace",
		"Name",
		"Type",
		"Replicas",
		"Rollout Strategy",
		"Anti-Affinity",
		"PVCs",
		"Rank",
	})

	for _, v := range p.Items {
		table.Append([]string{
			v.Namespace,
			v.Name,
			v.Type,
			fmt.Sprint(v.Replicas),
			v.RolloutStrategy,
			v.PodAntiAffinity,
			v.PVC,
			fmt.Sprint(v.Rank),
		})
	}

	table.Render()
	content := buf.String()
	return fmt.Sprintf("%s %s", content, description)
}
