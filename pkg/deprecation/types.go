package deprecation

import (
	"bytes"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
)

type Deprecation struct {
	KubeOverrides           *clientcmd.ConfigOverrides
	KubeClient              *kubernetes.KubeClient
	K8sVersion              string
	IgnoredKinds            []string
	ThrottleBurst           int
	Category                string
	Subject                 string
	ValidationFailureAction string
	Background              bool
}

type GroupVersion struct {
	Group     string     `json:"group"`
	Version   string     `json:"version"`
	Resources []Resource `json:"resources"`
}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type Resource struct {
	Name        string           `json:"name"`
	Introduced  string           `json:"introduced"`
	Deprecated  string           `json:"deprecated"`
	Removed     string           `json:"removed"`
	Replacement GroupVersionKind `json:"replacement"`
}

type DeprecatedResource struct {
	Namespace               string
	Name                    string
	DeprecatedGroupVersion  string
	DeprecatedKind          string
	DeprecationVersion      string
	RemovalVersion          string
	ReplacementGroupVersion string
	ReplacementKind         string
}

// DeprecatedResourceList implements PrintableList
type DeprecatedResourceList struct {
	Items []DeprecatedResource
}

func (p DeprecatedResourceList) ToJson() string {
	return printer.ToJson(p.Items)
}

func (p DeprecatedResourceList) ToHtml() string {
	return printer.ToHtml(HtmlPage, p)
}

func (p DeprecatedResourceList) ToTable() string {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{
		"Namespace",
		"Name",
		"Deprecated GV",
		"Deprecated Kind",
		"Replacement GV",
		"Replacement Kind",
		"Deprecation Version",
		"Removal Version",
	})

	for _, v := range p.Items {
		table.Append([]string{
			v.Namespace,
			v.Name,
			v.DeprecatedGroupVersion,
			v.DeprecatedKind,
			v.ReplacementGroupVersion,
			v.ReplacementKind,
			v.DeprecationVersion,
			v.RemovalVersion,
		})
	}

	table.Render()
	return buf.String()
}
