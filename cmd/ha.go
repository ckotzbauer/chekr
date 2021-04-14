package cmd

import (
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/ha"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/spf13/cobra"
)

// haCmd represents the ha command
var haCmd = &cobra.Command{
	Use:   "ha",
	Short: "Creates high-availability report of your workload.",
	Run: func(cmd *cobra.Command, args []string) {
		selector, _ := cmd.Flags().GetString("selector")
		namespace, _ := cmd.Flags().GetString("namespace")

		r := ha.HighAvailability{
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
	rootCmd.AddCommand(haCmd)
	haCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	// Output
}
