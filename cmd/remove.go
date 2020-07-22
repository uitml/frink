package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
)

type RemoveContext struct {
	cli.CommandContext

	WaitForDelete bool
}

func NewRemoveCmd() *cobra.Command {
	ctx := &RemoveContext{}
	cmd := &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove job from cluster",

		PreRunE: ctx.PreRun,
		RunE:    ctx.Run,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&ctx.WaitForDelete, "wait", "w", false, "wait for job to be deleted")

	return cmd
}

func (ctx *RemoveContext) PreRun(cmd *cobra.Command, args []string) error {
	return ctx.Initialize(cmd)
}

func (ctx *RemoveContext) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("job name must be specified")
	}

	name := args[0]

	job, err := ctx.Client.GetJob(name)
	if err != nil {
		return fmt.Errorf("unable to get job: %w", err)
	}

	if job == nil {
		fmt.Printf("Nothing to delete: no job named %s found\n", name)
		return nil
	}

	fmt.Printf("Deleting job %s...\n", name)
	if err := ctx.Client.DeleteJob(name); err != nil {
		return fmt.Errorf("unable to delete job: %w", err)
	}

	if !ctx.WaitForDelete {
		return nil
	}

	if err := waitUntilJobDeleted(name, ctx.Client); err != nil {
		return fmt.Errorf("timed out waiting for job to be deleted: %w", err)
	}

	return nil
}
