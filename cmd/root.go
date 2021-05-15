package cmd

import (
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	overrides *clientcmd.ConfigOverrides

	rootCmd = &cobra.Command{
		Use:   "chekr",
		Short: "A inspection utility for kubernetes-cluster maintenance.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")

			if output != "table" && output != "json" && output != "html" {
				return fmt.Errorf("Output-Format not valid: %v", output)
			}

			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output-Format. Valid values are [table, json, html]")
	rootCmd.PersistentFlags().StringP("output-file", "", "", "File to write to output to.")

	overrides = kubernetes.BindFlags(rootCmd.PersistentFlags())
}
