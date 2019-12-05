package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s"
)

// NOTE: Global package state; bad idea, but works for the time being.
var (
	kubectx *k8s.KubeContext
)

var rootCmd = &cobra.Command{
	Use:   "frink",
	Short: "Frink simplifies your Springfield workflows",

	// Silence usage when an error occurs.
	SilenceUsage: true,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := k8s.Client("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %w", err)
		}
		kubectx = ctx

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(runCmd)

	cli.DisableFlagsInUseLine(rootCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
