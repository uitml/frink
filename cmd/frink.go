package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(runCmd)

	disableFlagsInUseLine(rootCmd)
}

// visitAll visits the entire command tree rooted at cmd, invoking fn on each command.
func visitAll(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, child := range cmd.Commands() {
		visitAll(child, fn)
	}
}

// disableFlagsInUseLine sets the disableFlagsInUseLine flag on the entire command tree rooted at cmd.
func disableFlagsInUseLine(cmd *cobra.Command) {
	visitAll(cmd, func(c *cobra.Command) {
		c.DisableFlagsInUseLine = true
	})
}
