package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s"
)

type LogsContext struct {
	cli.CommandContext
}

func NewLogsCmd() *cobra.Command {
	ctx := &LogsContext{}
	cmd := &cobra.Command{
		Use:     "logs <name>",
		Short:   "Fetch the logs of a job",
		Aliases: []string{"watch"}, // TODO(thomasjo): Remove alias?

		PreRunE: ctx.PreRun,
		RunE:    ctx.Run,
	}

	return cmd
}

func (ctx *LogsContext) PreRun(cmd *cobra.Command, args []string) error {
	return ctx.Initialize(cmd)
}

func (ctx *LogsContext) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("job name must be specified")
	}

	name := args[0]
	req, err := ctx.Client.GetJobLogs(name, k8s.DefaultLogOptions)
	if err != nil {
		return fmt.Errorf("unable to get logs: %w", errors.Unwrap(err))
	}

	if req == nil {
		return fmt.Errorf("unable to get logs: request not returned (nil)")
	}

	stream, err := req.Stream()
	if err != nil {
		return fmt.Errorf("unable to stream logs: %w", err)
	}
	defer stream.Close()

	reader := bufio.NewReader(stream)
	if _, err := io.Copy(os.Stdout, reader); err != nil {
		return fmt.Errorf("unable to write output: %w", err)
	}

	return nil
}
