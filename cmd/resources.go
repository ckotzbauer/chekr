package cmd

import (
	"time"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/resources"
	"github.com/prometheus/common/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Analyze resource requests and limits of pods.",
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("prometheus-url")
		username := viper.GetString("prometheus-username")
		password := viper.GetString("prometheus-password")
		countDays := viper.GetInt64("count-days")
		timeout := viper.GetDuration("timeout")

		labelSelector := viper.GetString("selector")
		annotationSelector := viper.GetString("annotation")
		namespace := viper.GetString("namespace")
		cpuMetric := viper.GetString("cpu-metric")
		memoryMetric := viper.GetString("memory-metric")
		limitsThreshold := viper.GetInt("limits-threshold")
		requestsThreshold := viper.GetInt("requests-threshold")

		r := resources.Resource{
			Prometheus: prometheus.Prometheus{
				Url:       url,
				UserName:  username,
				Password:  config.Secret(password),
				CountDays: countDays,
				Timeout:   timeout,
			},
			KubeClient:         kubernetes.NewClient(cmd, overrides),
			Pods:               args,
			LabelSelector:      labelSelector,
			AnnotationSelector: annotationSelector,
			Namespace:          namespace,
			CpuMetric:          cpuMetric,
			MemoryMetric:       memoryMetric,
			LimitsThreshold:    limitsThreshold,
			RequestsThreshold:  requestsThreshold,
		}

		list := r.Execute()

		output := viper.GetString("output")
		outputFile := viper.GetString("output-file")

		printer := printer.Printer{Type: output, File: outputFile}
		printer.Print(list)
	},
}

func init() {
	rootCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().StringP("prometheus-url", "u", "", "Prometheus-URL")
	resourcesCmd.Flags().StringP("prometheus-username", "U", "", "Prometheus-Username")
	resourcesCmd.Flags().StringP("prometheus-password", "P", "", "Prometheus-Password")
	resourcesCmd.Flags().String("cpu-metric", "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate", "CPU-Usage metric to query")
	resourcesCmd.Flags().String("memory-metric", "container_memory_working_set_bytes", "Memory-Usage metric to query")
	resourcesCmd.Flags().Int64P("count-days", "d", 30, "Count of days to analyze metrics from (until now).")
	resourcesCmd.Flags().DurationP("timeout", "t", time.Duration(30)*time.Second, "Timeout")
	resourcesCmd.Flags().Int("limits-threshold", -1, "Only emit pods with a greater deviation of applied limits in average.")
	resourcesCmd.Flags().Int("requests-threshold", -1, "Only emit pods with a greater deviation of applied requests in average.")
	resourcesCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	resourcesCmd.Flags().StringP("annotation", "a", "", "Annotation-Selector")
	resourcesCmd.MarkFlagRequired("prometheus-url")

	viper.BindPFlag("prometheus-url", resourcesCmd.Flags().Lookup("prometheus-url"))
	viper.BindPFlag("prometheus-username", resourcesCmd.Flags().Lookup("prometheus-username"))
	viper.BindPFlag("prometheus-password", resourcesCmd.Flags().Lookup("prometheus-password"))
	viper.BindPFlag("cpu-metric", resourcesCmd.Flags().Lookup("cpu-metric"))
	viper.BindPFlag("memory-metric", resourcesCmd.Flags().Lookup("memory-metric"))
	viper.BindPFlag("count-days", resourcesCmd.Flags().Lookup("count-days"))
	viper.BindPFlag("timeout", resourcesCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("limits-threshold", resourcesCmd.Flags().Lookup("limits-threshold"))
	viper.BindPFlag("requests-threshold", resourcesCmd.Flags().Lookup("requests-threshold"))
	viper.BindPFlag("selector", resourcesCmd.Flags().Lookup("selector"))
	viper.BindPFlag("annotation", resourcesCmd.Flags().Lookup("annotation"))
}
