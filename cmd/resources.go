package cmd

import (
	"fmt"
	"time"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/prometheus"
	"github.com/ckotzbauer/chekr/pkg/resources"
	"github.com/prometheus/common/config"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var overrides *clientcmd.ConfigOverrides

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Analyze resource requests and limits of pods.",
	Args: func(cmd *cobra.Command, args []string) error {
		labelSelector, _ := cmd.Flags().GetString("selector")
		if labelSelector == "" && len(args) < 1 {
			return fmt.Errorf("Pod argument missing. Otherwise use selector flag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("prometheus-url")
		username, _ := cmd.Flags().GetString("prometheus-username")
		password, _ := cmd.Flags().GetString("prometheus-password")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		selector, _ := cmd.Flags().GetString("selector")
		namespace, _ := cmd.Flags().GetString("namespace")

		r := resources.Resource{
			Prometheus: prometheus.Prometheus{
				Url:      url,
				UserName: username,
				Password: config.Secret(password),
				Timeout:  timeout,
			},
			KubeOverrides: overrides,
			KubeClient:    kubernetes.NewClient(overrides),
			Pods:          args,
			Selector:      selector,
			Namespace:     namespace,
		}

		list, err := r.Execute()

		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

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
	resourcesCmd.Flags().DurationP("timeout", "t", time.Duration(30)*time.Second, "Timeout")

	resourcesCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	// Output

	overrides = kubernetes.BindFlags(resourcesCmd.Root().PersistentFlags())
}
