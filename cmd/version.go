package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
)

// Build metadata that is typically overriden by build tools.
var (
	version = "0.0.0-dev"
	commit  = "unknown"
	date    = time.Now().Format(time.RFC3339)
)

type VersionContext struct {
	cli.CommandContext
}

func NewVersionCmd() *cobra.Command {
	ctx := &VersionContext{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version information",

		Run: ctx.Run,
	}

	return cmd
}

func (ctx *VersionContext) Run(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "frink %s\n", version)
}
