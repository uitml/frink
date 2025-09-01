// Package cmd provides implementations of CLI commands.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frink",
		Short: "Frink simplifies your Springfield workflows",

		// Do not display usage when an error occurs.
		SilenceUsage: true,
	}

	pflags := cmd.PersistentFlags()
	pflags.String("context", "", "name of the kubeconfig context to use")
	pflags.StringP("namespace", "n", "", "cluster namespace to use")

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newLogsCmd())
	cmd.AddCommand(newRemoveCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newDebugCmd())
	cmd.AddCommand(newGPUCmd())
	cli.DisableFlagsInUseLine(cmd)

	return cmd
}

// Execute executes the root command using os.Args, running through the command tree and invoking the matching subcommand.
func Execute() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
