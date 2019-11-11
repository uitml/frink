package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Schedule a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job specification file must be specified")
		}

		file := args[0]
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("specified file does not exist: %v", file)
			}

			return fmt.Errorf("unable to access file: %v", err)
		}

		job, err := k8s.ParseJob(file)
		if err != nil {
			return fmt.Errorf("unable to parse job: %v", err)
		}

		// TODO: Reconsider this? Many reasons to avoid this; should be challenged.
		k8s.OverrideJobSpec(job)

		kubectx, err := k8s.Client("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %v", err)
		}

		err = kubectx.DeleteJob(job.Name)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("unable to previous job: %v", err)
		}

		// Try to create the job using retry with backoff.
		// This handles scenarios where an existing job is still being terminated, etc.
		err = k8s.RetryOnExists(k8s.DefaultBackoff, func() error { return kubectx.CreateJob(job) })
		if err != nil {
			return fmt.Errorf("unable to create job: %v", err)
		}

		// TODO: Implement support for streaming job/pod log to stdout.

		return nil
	},
}
