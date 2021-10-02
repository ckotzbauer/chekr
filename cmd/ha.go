package cmd

import (
	"github.com/ckotzbauer/chekr/pkg/ha"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// haCmd represents the ha command
var haCmd = &cobra.Command{
	Use:   "ha",
	Short: "Creates high-availability report of your workload.",
	Run: func(cmd *cobra.Command, args []string) {
		labelSelector := viper.GetString("selector")
		annotationSelector := viper.GetString("annotation")
		namespace := viper.GetString("namespace")

		r := ha.HighAvailability{
			KubeClient:         kubernetes.NewClient(cmd, overrides),
			Pods:               args,
			LabelSelector:      labelSelector,
			AnnotationSelector: annotationSelector,
			Namespace:          namespace,
		}

		list := r.Execute()

		output := viper.GetString("output")
		outputFile := viper.GetString("output-file")

		printer := printer.Printer{Type: output, File: outputFile}
		printer.Print(list)
	},
}

func init() {
	rootCmd.AddCommand(haCmd)
	haCmd.Flags().StringP("selector", "l", "", "Label-Selector")
	haCmd.Flags().StringP("annotation", "a", "", "Annotation-Selector")

	viper.BindPFlag("selector", haCmd.Flags().Lookup("selector"))
	viper.BindPFlag("annotation", haCmd.Flags().Lookup("annotation"))
}
