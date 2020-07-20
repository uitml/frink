package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
	"github.com/uitml/frink/internal/k8s/retry"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

var backoff = wait.Backoff{
	Duration: 100 * time.Millisecond,
	Factor:   1.0,
	Steps:    1200,
}

var follow bool

var runCmd = &cobra.Command{
	Use:   "run <file>",
	Short: "Schedule a job on the cluster",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job specification file must be specified")
		}

		file := args[0]
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("specified file does not exist: %v", file)
			}

			return fmt.Errorf("unable to access file: %w", err)
		}

		job, err := k8s.ParseJob(file)
		if err != nil {
			return fmt.Errorf("unable to parse job: %w", err)
		}

		// TODO: Reconsider this? Many reasons to avoid this; should be challenged.
		k8s.OverrideJobSpec(job)

		if err := deletePreviousJob(job); err != nil {
			return fmt.Errorf("unable to delete previous job: %w", err)
		}

		// Try to create the job using retry.
		// This handles scenarios where an existing job is still being terminated, etc.
		fmt.Println("Creating job...")
		err = retry.OnExists(backoff, func() error { return kubectx.CreateJob(job) })
		if err != nil {
			return fmt.Errorf("unable to create job: %w", err)
		}

		if !follow {
			return nil
		}

		// TODO: Ensure nil references are properly handled in this block.
		err = retry.OnError(backoff, apierrors.IsBadRequest, func() error {
			req, err := kubectx.GetJobLogs(job.Name, k8s.DefaultLogOptions)
			if err != nil {
				return errors.Unwrap(err)
			}

			if req == nil {
				return fmt.Errorf("unable to get logs: request not returned (nil)")
			}

			stream, err := req.Stream()
			if err != nil {
				return err
			}
			defer stream.Close()

			reader := bufio.NewReader(stream)
			if _, err := io.Copy(os.Stdout, reader); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.BoolVarP(&follow, "follow", "f", false, "wait for job to start, then stream logs")
}

func deletePreviousJob(job *batchv1.Job) error {
	oldJob, err := kubectx.GetJob(job.Name)
	if err != nil {
		return fmt.Errorf("unable to get previous job: %w", err)
	}

	if oldJob != nil {
		fmt.Println("Deleting previous job...")
		err = kubectx.DeleteJob(oldJob.Name)
		if err != nil {
			return err
		}

		if err := waitUntilDeleted(oldJob); err != nil {
			return err
		}
	}

	return nil
}

func waitUntilDeleted(job *batchv1.Job) error {
	err := wait.Poll(100*time.Millisecond, 120*time.Second, func() (bool, error) {
		oldJob, err := kubectx.GetJob(job.Name)
		if err != nil {
			return false, err
		}

		return oldJob == nil, nil
	})

	return err
}
