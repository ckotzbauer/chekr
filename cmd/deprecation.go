package cmd

import (
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/deprecation"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/spf13/cobra"
)

// deprecationCmd represents the deprecation command
var deprecationCmd = &cobra.Command{
	Use:   "deprecation",
	Short: "Lists deprecated objects in cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		ignoredKinds, _ := cmd.Flags().GetStringSlice("ignored-kinds")

		r := deprecation.Deprecation{
			KubeOverrides: overrides,
			KubeClient:    kubernetes.NewClient(overrides),
			Namespace:     namespace,
			K8sVersion:    k8sVersion,
			IgnoredKinds:  ignoredKinds,
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
	rootCmd.AddCommand(deprecationCmd)
	deprecationCmd.Flags().StringP("k8s-version", "V", "", "Highest K8s major.minor version to show deprecations for (e.g. 1.21)")
	deprecationCmd.Flags().StringSliceP("ignored-kinds", "i", []string{}, "All kinds you want to ignore (e.g. Deployment,DaemonSet)")
	// Output
}
