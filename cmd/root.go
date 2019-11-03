package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "frink",
	Short: "Simplifies your Springfield workflows",

	// Silence usage when an error occurs.
	SilenceUsage: true,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(jobCmd)

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
