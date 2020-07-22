// Package cmd provides implementations of CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s"
)

// NOTE: Global package state; bad idea, but works for the time being.
var client k8s.KubeClient

var rootCmd = &cobra.Command{
	Use:   "frink",
	Short: "Frink simplifies your Springfield workflows",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		context := viper.GetString("context")
		namespace := viper.GetString("namespace")

		kubectx, err := k8s.Client(context, namespace)
		if err != nil {
			return fmt.Errorf("unable to get kube client: %w", err)
		}

		client = kubectx
		return nil
	},

	// Do not display usage when an error occurs.
	SilenceUsage: true,
}

func init() {
	cobra.OnInitialize(cli.InitConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.String("context", "", "name of the kubeconfig context to use")
	pflags.StringP("namespace", "n", "", "cluster namespace to use")
	viper.BindPFlags(pflags)

	rootCmd.AddCommand(NewListCmd())

	cli.DisableFlagsInUseLine(rootCmd)
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
