package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const completionDesc = `
Generate autocompletion scripts for Chekr for the specified shell.
`
const bashCompDesc = `
Generate the autocompletion script for Chekr for the bash shell.
To load completions in your current shell session:
$ source <(chekr completion bash)
To load completions for every new session, execute once:
Linux:
  $ chekr completion bash > /etc/bash_completion.d/chekr
MacOS:
  $ chekr completion bash > /usr/local/etc/bash_completion.d/chekr
`

const zshCompDesc = `
Generate the autocompletion script for Chekr for the zsh shell.
To load completions in your current shell session:
$ source <(chekr completion zsh)
To load completions for every new session, execute once:
$ chekr completion zsh > "${fpath[1]}/_chekr"
`

const fishCompDesc = `
Generate the autocompletion script for Chekr for the fish shell.
To load completions in your current shell session:
$ chekr completion fish | source
To load completions for every new session, execute once:
$ chekr completion fish > ~/.config/fish/completions/chekr.fish
You will need to start a new shell for this setup to take effect.
`

const powershellCompDesc = `
Generate the autocompletion script for Chekr for the powershell.
To load completions in your current shell session:
PS> chekr completion powershell | Out-String | Invoke-Expression
To load completions for every new session, run:
PS> chekr completion powershell > yourprogram.ps1
# and source this file from your PowerShell profile.
`

const (
	noDescFlagName = "no-descriptions"
	noDescFlagText = "disable completion descriptions"
)

var disableCompDescriptions bool

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "generate autocompletion scripts for the specified shell",
		Long:  completionDesc,
		Args:  noArgs,
	}

	bash := &cobra.Command{
		Use:                   "bash",
		Short:                 "generate autocompletion script for bash",
		Long:                  bashCompDesc,
		Args:                  noArgs,
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletion(os.Stdout)
		},
	}

	zsh := &cobra.Command{
		Use:               "zsh",
		Short:             "generate autocompletion script for zsh",
		Long:              zshCompDesc,
		Args:              noArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if disableCompDescriptions {
				return rootCmd.GenZshCompletionNoDesc(os.Stdout)
			} else {
				return rootCmd.GenZshCompletion(os.Stdout)
			}
		},
	}
	zsh.Flags().BoolVar(&disableCompDescriptions, noDescFlagName, false, noDescFlagText)

	fish := &cobra.Command{
		Use:               "fish",
		Short:             "generate autocompletion script for fish",
		Long:              fishCompDesc,
		Args:              noArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(os.Stdout, disableCompDescriptions)
		},
	}
	fish.Flags().BoolVar(&disableCompDescriptions, noDescFlagName, false, noDescFlagText)

	powershell := &cobra.Command{
		Use:               "powershell",
		Short:             "generate autocompletion script for powershell",
		Long:              powershellCompDesc,
		Args:              noArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if disableCompDescriptions {
				return rootCmd.GenPowerShellCompletion(os.Stdout)
			} else {
				return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	powershell.Flags().BoolVar(&disableCompDescriptions, noDescFlagName, false, noDescFlagText)

	cmd.AddCommand(bash, zsh, fish, powershell)

	return cmd
}

// NoArgs returns an error if any args are included.
func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf(
			"%q accepts no arguments\n\nUsage:  %s",
			cmd.CommandPath(),
			cmd.UseLine(),
		)
	}
	return nil
}

// Function to disable file completion
func noCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(newCompletionCmd())
}
