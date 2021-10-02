package deprecation

import (
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/deprecation"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

func createKyvernoCreateCmd(overrides *clientcmd.ConfigOverrides) *cobra.Command {
	// kyvernoCreateCmd represents the deprecation kyverno-create command
	kyvernoCreateCmd := &cobra.Command{
		Use:   "kyverno-create",
		Short: "Creates Kyverno validation policies for deprecated objects in cluster.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			output := viper.GetString("output")

			if output != "yaml" && output != "json" {
				return fmt.Errorf("Output-Format not valid: %v", output)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			k8sVersion := viper.GetString("k8s-version")
			ignoredKinds := viper.GetStringSlice("ignored-kinds")
			category := viper.GetString("category")
			subject := viper.GetString("subject")
			validationFailureAction := viper.GetString("validation-failure-action")
			background := viper.GetBool("background")

			r := deprecation.Deprecation{
				KubeClient:              kubernetes.NewClient(cmd, overrides),
				K8sVersion:              k8sVersion,
				IgnoredKinds:            ignoredKinds,
				Category:                category,
				Subject:                 subject,
				ValidationFailureAction: validationFailureAction,
				Background:              background,
			}

			policyCrd := r.ExecuteKyvernoCreate()
			stringOutput := policyCrd

			output := viper.GetString("output")
			outputFile := viper.GetString("output-file")
			dryRun := viper.GetBool("dry-run")

			r.HandleKyvernoResult(stringOutput, output, outputFile, dryRun)
		},
	}

	kyvernoCreateCmd.PersistentFlags().StringP("output", "o", "yaml", "Output-Format. Valid values are [yaml, json]")
	kyvernoCreateCmd.Flags().StringP("k8s-version", "V", "", "Highest K8s major.minor version to show deprecations for (e.g. 1.21)")
	kyvernoCreateCmd.Flags().StringSliceP("ignored-kinds", "i", []string{}, "All kinds you want to ignore (e.g. Deployment,DaemonSet)")
	kyvernoCreateCmd.Flags().String("category", "Best Practices", "Category set for 'policies.kyverno.io/category' annotation.")
	kyvernoCreateCmd.Flags().String("subject", "Kubernetes APIs", "Subject set for 'policies.kyverno.io/subject' annotation.")
	kyvernoCreateCmd.Flags().String("validation-failure-action", "audit", "Validation-Failure-Action of the policy (audit or failure).")
	kyvernoCreateCmd.Flags().Bool("background", true, "Whether background scans should be performed.")
	kyvernoCreateCmd.Flags().Bool("dry-run", false, "Whether or not the generated policy should be applied.")

	viper.BindPFlag("k8s-version", kyvernoCreateCmd.Flags().Lookup("k8s-version"))
	viper.BindPFlag("ignored-kinds", kyvernoCreateCmd.Flags().Lookup("ignored-kinds"))
	viper.BindPFlag("category", kyvernoCreateCmd.Flags().Lookup("category"))
	viper.BindPFlag("k8s-subject", kyvernoCreateCmd.Flags().Lookup("subject"))
	viper.BindPFlag("validation-failure-action", kyvernoCreateCmd.Flags().Lookup("validation-failure-action"))
	viper.BindPFlag("background", kyvernoCreateCmd.Flags().Lookup("background"))
	viper.BindPFlag("dry-run", kyvernoCreateCmd.Flags().Lookup("dry-run"))
	return kyvernoCreateCmd
}

func InitKyvernoCreateCmd(deprecationCmd *cobra.Command, overrides *clientcmd.ConfigOverrides) {
	deprecationCmd.AddCommand(createKyvernoCreateCmd(overrides))
	// Output
}
