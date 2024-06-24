package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
)

type debugContext struct {
	cli.CommandContext
	Namespace    string
	ResourceName string
	ResourceType string // "pod" or "job"
}

func newDebugCmd() *cobra.Command {
	ctx := &debugContext{}
	cmd := &cobra.Command{
		Use:     "debug [resource-name]",
		Short:   "Debug resources by showing events",
		Args:    cobra.ExactArgs(1), // Ensures exactly one argument is passed
		PreRunE: ctx.PreRun,
		RunE:    ctx.Run,
	}

	flags := cmd.Flags()
	flags.StringVarP(&ctx.Namespace, "namespace", "n", "default", "Specify the namespace of the resource")

	return cmd
}

func (ctx *debugContext) PreRun(cmd *cobra.Command, args []string) error {
	ctx.ResourceName = args[0]
	return ctx.Initialize(cmd)
}

func (ctx *debugContext) Run(cmd *cobra.Command, args []string) error {
	// Try to get job events and associated pod events
	jobEvents, err := ctx.Client.GetJobEvents(ctx.ResourceName)
	if err == nil && jobEvents != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "Job Events:\n", jobEvents)
		podNames, err := ctx.Client.GetPodsFromJob(ctx.ResourceName)
		if err == nil {
			for _, podName := range podNames {
				podEvents, _ := ctx.Client.GetPodEvents(podName)
				if podEvents != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "Pod Events for %s:\n%s\n", podName, podEvents)
				}
			}
		}
		return nil
	}

	// Try to get pod events and associated job events
	podEvents, err := ctx.Client.GetPodEvents(ctx.ResourceName)
	if err == nil && podEvents != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "Pod Events:\n", podEvents)
		jobName, err := ctx.Client.GetJobFromPod(ctx.ResourceName)
		if err == nil && jobName != "" {
			jobEvents, _ := ctx.Client.GetJobEvents(jobName)
			if jobEvents != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Job Events for %s:\n%s\n", jobName, jobEvents)
			}
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not get events: %w", err)
	}

	return fmt.Errorf("no events found for %s in namespace %s", ctx.ResourceName, ctx.Namespace)
}
