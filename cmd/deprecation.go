package cmd

import (
	"github.com/ckotzbauer/chekr/cmd/deprecation"
	"github.com/spf13/cobra"
)

// deprecationCmd represents the deprecation command
var deprecationCmd = &cobra.Command{
	Use:   "deprecation",
	Short: "Handle deprecated objects in cluster.",
}

func init() {
	rootCmd.AddCommand(deprecationCmd)

	deprecation.InitListCmd(deprecationCmd, overrides)
	deprecation.InitKyvernoCreateCmd(deprecationCmd, overrides)
	deprecation.InitKyvernoDeleteCmd(deprecationCmd, overrides)
	// Output
}
