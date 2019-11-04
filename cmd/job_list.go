package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
	batchv1 "k8s.io/api/batch/v1"
)

var (
	headerNames = []string{
		"NAME",
		"SUCCEEDED",
	}
)

// Flags
var (
	showAll bool
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubectx, err := k8s.Client("")
		if err != nil {
			return fmt.Errorf("could not get k8s client: %v", err)
		}

		jobs, err := kubectx.ListJobs()
		if err != nil {
			return fmt.Errorf("could not list jobs: %v", err)
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

func init() {
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show all jobs (defaults to only active)")
}

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
	data := []string{
		job.Name,
		strconv.FormatInt(int64(job.Status.Succeeded), 10),
	}

	return strings.Join(data, "\t") + "\t" // TODO: Check if the trailing tab is required
}
