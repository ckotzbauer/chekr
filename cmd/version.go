package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func NewVersionCmd(version, commit, date, builtBy string) *cobra.Command {
	// versionCmd represents the version command
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of chekr.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %v \n", version)
			fmt.Printf("Commit: %v \n", commit)
			fmt.Printf("Buit at: %v \n", date)
			fmt.Printf("Buit by: %v \n", builtBy)
			fmt.Printf("Go Version: %v \n", runtime.Version())
		},
	}
}
