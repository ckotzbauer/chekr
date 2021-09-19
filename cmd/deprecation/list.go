package deprecation

import (
	"os"

	"github.com/ckotzbauer/chekr/pkg/deprecation"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func createListCmd(overrides *clientcmd.ConfigOverrides) *cobra.Command {
	// listCmd represents the deprecation list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists deprecated objects in cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			k8sVersion, _ := cmd.Flags().GetString("k8s-version")
			ignoredKinds, _ := cmd.Flags().GetStringSlice("ignored-kinds")
			throttleBurst, _ := cmd.Flags().GetInt("throttle-burst")

			r := deprecation.Deprecation{
				KubeOverrides: overrides,
				KubeClient:    kubernetes.NewClient(cmd, overrides),
				K8sVersion:    k8sVersion,
				IgnoredKinds:  ignoredKinds,
				ThrottleBurst: throttleBurst,
			}

			list := r.Execute()

			output, _ := cmd.Flags().GetString("output")
			outputFile, _ := cmd.Flags().GetString("output-file")
			omitExitCode, _ := cmd.Flags().GetBool("omit-exit-code")

			printer := printer.Printer{Type: output, File: outputFile}
			printer.Print(list)

			items := list.(deprecation.DeprecatedResourceList)
			if len(items.Items) > 0 && !omitExitCode {
				os.Exit(1)
			}
		},
	}

	listCmd.Flags().StringP("k8s-version", "V", "", "Highest K8s major.minor version to show deprecations for (e.g. 1.21)")
	listCmd.Flags().StringSliceP("ignored-kinds", "i", []string{}, "All kinds you want to ignore (e.g. Deployment,DaemonSet)")
	listCmd.Flags().Bool("omit-exit-code", false, "Omits the non-zero exit code if deprecations were found.")
	listCmd.Flags().IntP("throttle-burst", "t", 100, "Burst used for throttling of Kubernetes discovery-client")
	return listCmd
}

func InitListCmd(deprecationCmd *cobra.Command, overrides *clientcmd.ConfigOverrides) {
	deprecationCmd.AddCommand(createListCmd(overrides))
	// Output
}
