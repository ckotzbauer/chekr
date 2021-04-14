package ha

import (
	"bytes"
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
)

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

func (p PodAvailabilityList) ToJson() (string, error) {
	return printer.ToJson(p.Items)
}

func (p PodAvailabilityList) ToHtml() (string, error) {
	return printer.ToHtml(HtmlPage, p)
}

func (p PodAvailabilityList) ToTable() (string, error) {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{
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
	return buf.String(), nil
}
