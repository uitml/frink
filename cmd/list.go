package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
)

var (
	showAll bool

	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List jobs",

		RunE: func(cmd *cobra.Command, args []string) error {
			jobs, err := kubectx.ListJobs()
			if err != nil {
				return fmt.Errorf("could not list jobs: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			defer w.Flush()

			fmt.Fprintln(w, header())
			for _, job := range jobs.Items {
				fmt.Fprintln(w, row(job))
			}

			return nil
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show all jobs; active and terminated")
}

var (
	headerNames = []string{
		"NAME",
		"STATUS",
		"DURATION",
		"AGE",
	}
)

func header() string {
	return strings.Join(headerNames, "\t") + "\t" // TODO: Check if the trailing tab is required
}

func underlinedHeader() string {
	var rules []string
	for _, name := range headerNames {
		rules = append(rules, strings.Repeat("-", len(name)))
	}

	head := header()
	rule := strings.Join(rules, "\t")

	return fmt.Sprintf("%s\n%s\t", head, rule)
}

func row(job batchv1.Job) string {
	total := job.Status.Active + job.Status.Succeeded + job.Status.Failed
	succeeded := job.Status.Succeeded

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

	data := []string{
		job.Name,
		fmt.Sprintf("%d/%d", succeeded, total),
		durafmt.Parse(duration).LimitFirstN(2).String(),
		humanize.Time(start),
	}

	return strings.Join(data, "\t") + "\t" // TODO: Check if the trailing tab is required
}
