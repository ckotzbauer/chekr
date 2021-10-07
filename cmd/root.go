package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	overrides *clientcmd.ConfigOverrides
	cfgFile   string
	verbosity string

	rootCmd = &cobra.Command{
		Use:   "chekr",
		Short: "A inspection utility for kubernetes-cluster maintenance.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			output := viper.GetString("output")

			if output != "table" && output != "json" && output != "html" {
				return fmt.Errorf("Output-Format not valid: %v", output)
			}

			return nil
		},
	}
)

// Execute executes the root command.
func Execute(version, commit, date, builtBy string) error {
	rootCmd.AddCommand(NewVersionCmd(version, commit, date, builtBy))
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the chekr config-file.")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output-Format. Valid values are [table, json, html]")
	rootCmd.PersistentFlags().StringP("output-file", "", "", "File to write to output to.")
	rootCmd.PersistentFlags().String(clientcmd.RecommendedConfigPathFlag, "", "Path to the kubeconfig file to use for CLI requests.")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", logrus.WarnLevel.String(), "Log-level (debug, info, warn, error, fatal, panic)")

	overrides = kubernetes.BindFlags(rootCmd.PersistentFlags())

	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("output-file", rootCmd.PersistentFlags().Lookup("output-file"))
	viper.BindPFlag(clientcmd.RecommendedConfigPathFlag, rootCmd.PersistentFlags().Lookup(clientcmd.RecommendedConfigPathFlag))
	viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	if err := setUpLogs(os.Stdout, verbosity); err != nil {
		fmt.Println(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("chekr")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/chekr")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("CHEKR")

	if err := viper.ReadInConfig(); err != nil && cfgFile != "" {
		logrus.WithError(err).Fatalf("An error occurred while reading the config!")
	}
}

//setUpLogs set the log output ans the log level
func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)
	return nil
}
