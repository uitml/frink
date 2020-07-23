package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
	batchv1 "k8s.io/api/batch/v1"
)

type listContext struct {
	cli.CommandContext

	ShowAll bool
}

func newListCmd() *cobra.Command {
	ctx := &listContext{}
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List jobs",

		PreRunE: ctx.PreRun,
		RunE:    ctx.Run,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&ctx.ShowAll, "all", "a", false, "show all jobs; active and terminated")

	return cmd
}

func (ctx *listContext) PreRun(cmd *cobra.Command, args []string) error {
	return ctx.Initialize(cmd)
}

func (ctx *listContext) Run(cmd *cobra.Command, args []string) error {
	jobs, err := ctx.Client.ListJobs()
	if err != nil {
		return fmt.Errorf("could not list jobs: %w", err)
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, header())
	for _, job := range jobs.Items {
		fmt.Fprintln(w, row(job))
	}

	return nil
}

func header() string {
	columnNames := []string{
		"NAME",
		"STATUS",
		"COMPLETIONS",
		"DURATION",
		"AGE",
	}

	return strings.Join(columnNames, "\t") + "\t"
}

func row(job batchv1.Job) string {
	columns := []string{
		job.Name,
		status(job),
		completions(job),
		duration(job),
		age(job),
	}

	return strings.Join(columns, "\t") + "\t"
}

func status(job batchv1.Job) string {
	switch {
	case job.Status.Active > 0:
		return "Active"
	case job.Spec.Completions == nil || *job.Spec.Completions == job.Status.Succeeded:
		return "Succeeded"
	case job.Status.Failed > 0:
		return "Failed"
	}

	return "Stopped"
}

func completions(job batchv1.Job) string {
	succeeded := job.Status.Succeeded
	total := succeeded + job.Status.Active + job.Status.Failed

	return fmt.Sprintf("%d/%d", succeeded, total)
}

func duration(job batchv1.Job) string {
	_, duration := timing(job)
	// TODO: Implement "smart" truncation scheme.
	humanized := durafmt.Parse(duration).LimitFirstN(2).String()

	return humanized
}

func age(job batchv1.Job) string {
	start, _ := timing(job)
	humanized := humanize.Time(start)

	return humanized
}

func timing(job batchv1.Job) (time.Time, time.Duration) {
	var start time.Time
	duration := time.Duration(0)
	if job.Status.StartTime != nil {
		start = job.Status.StartTime.Time
		end := time.Now()
		if job.Status.CompletionTime != nil {
			end = job.Status.CompletionTime.Time
		}
		duration = end.Sub(start)
	}

	return start, duration
}
