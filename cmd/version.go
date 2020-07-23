package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Build metadata that is typically overriden by build tools.
var (
	version = "0.0.0-dev"
	commit  = "unknown"
	date    = time.Now().Format(time.RFC3339)
)

type versionContext struct {
}

func newVersionCmd() *cobra.Command {
	ctx := &versionContext{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",

		Run: ctx.Run,
	}

	return cmd
}

func (ctx *versionContext) Run(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "frink %s\n", version)
}
