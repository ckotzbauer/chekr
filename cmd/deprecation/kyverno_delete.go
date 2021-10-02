package deprecation

import (
	"github.com/ckotzbauer/chekr/pkg/deprecation"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func createKyvernoDeleteCmd(overrides *clientcmd.ConfigOverrides) *cobra.Command {
	// kyvernoDeleteCmd represents the deprecation kyverno-delete command
	kyvernoDeleteCmd := &cobra.Command{
		Use:   "kyverno-delete",
		Short: "Deletes Kyverno validation policies for deprecated objects in cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			r := deprecation.Deprecation{
				KubeClient: kubernetes.NewClient(cmd, overrides),
			}

			r.DeletePolicy()
		},
	}

	return kyvernoDeleteCmd
}

func InitKyvernoDeleteCmd(deprecationCmd *cobra.Command, overrides *clientcmd.ConfigOverrides) {
	deprecationCmd.AddCommand(createKyvernoDeleteCmd(overrides))
	// Output
}
