package cmd

import (
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
		labelSelector, _ := cmd.Flags().GetString("selector")
		annotationSelector, _ := cmd.Flags().GetString("annotation")
		namespace, _ := cmd.Flags().GetString("namespace")

		r := ha.HighAvailability{
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
	rootCmd.AddCommand(haCmd)
	haCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	haCmd.Flags().StringP("annotation", "a", "", "Annotation-Selector")
	// Output
}
