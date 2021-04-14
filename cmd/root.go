package cmd

import (
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	cfgFile   string
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
	cobra.OnInitialize(initConfig)

	/*rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")*/

	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output-Format. Valid values are [table, json, html]")
	rootCmd.PersistentFlags().StringP("output-file", "", "", "File to write to output to.")

	overrides = kubernetes.BindFlags(rootCmd.PersistentFlags())

	//rootCmd.AddCommand(addCmd)
	//rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := homedir.Dir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		//viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
