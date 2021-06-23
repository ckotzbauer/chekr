package cmd

import (
	"time"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/resources"
	"github.com/prometheus/common/config"
	"github.com/spf13/cobra"
)

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Analyze resource requests and limits of pods.",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("prometheus-url")
		username, _ := cmd.Flags().GetString("prometheus-username")
		password, _ := cmd.Flags().GetString("prometheus-password")
		countDays, _ := cmd.Flags().GetInt64("count-days")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		labelSelector, _ := cmd.Flags().GetString("selector")
		annotationSelector, _ := cmd.Flags().GetString("annotation")
		namespace, _ := cmd.Flags().GetString("namespace")

		r := resources.Resource{
			Prometheus: prometheus.Prometheus{
				Url:       url,
				UserName:  username,
				Password:  config.Secret(password),
				CountDays: countDays,
				Timeout:   timeout,
			},
			KubeOverrides:      overrides,
			KubeClient:         kubernetes.NewClient(cmd, overrides),
			Pods:               args,
			LabelSelector:      labelSelector,
			AnnotationSelector: annotationSelector,
			Namespace:          namespace,
		}

		list := r.Execute()

		output, _ := cmd.Flags().GetString("output")
		outputFile, _ := cmd.Flags().GetString("output-file")

		printer := printer.Printer{Type: output, File: outputFile}
		printer.Print(list)
	},
}

func init() {
	rootCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().StringP("prometheus-url", "u", "", "Prometheus-URL")
	resourcesCmd.Flags().StringP("prometheus-username", "U", "", "Prometheus-Username")
	resourcesCmd.Flags().StringP("prometheus-password", "P", "", "Prometheus-Password")
	resourcesCmd.Flags().Int64P("count-days", "d", 30, "Count of days to analyze metrics from (until now).")
	resourcesCmd.Flags().DurationP("timeout", "t", time.Duration(30)*time.Second, "Timeout")

	resourcesCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	resourcesCmd.Flags().StringP("annotation", "a", "", "Annotation-Selector")

	resourcesCmd.MarkFlagRequired("prometheus-url")
}
