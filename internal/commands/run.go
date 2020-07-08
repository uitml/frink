package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	follow bool

	runCmd = &cobra.Command{
		Use:   "run [file]",
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

			err = kubectx.DeleteJob(job.Name)
			if err != nil && !apierrors.IsNotFound(err) {
				return fmt.Errorf("unable to previous job: %w", err)
			}

			// Try to create the job using retry with backoff.
			// This handles scenarios where an existing job is still being terminated, etc.
			err = k8s.RetryOnExists(k8s.DefaultBackoff, func() error { return kubectx.CreateJob(job) })
			if err != nil {
				return fmt.Errorf("unable to create job: %w", err)
			}

			if !follow {
				return nil
			}

			backoff := wait.Backoff{
				Duration: 1 * time.Second,
				Factor:   1.0,
				Steps:    120,
			}

			// TODO: Ensure nil references are properly handled in this block.
			err = k8s.OnError(backoff, apierrors.IsBadRequest, func() error {
				req, err := kubectx.GetJobLogs(job.Name, k8s.DefaultLogOptions)
				if err != nil {
					return errors.Unwrap(err)
				}

				if req == nil {
					// TODO: Inform user we did not get any logs?
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
)

func init() {
	runCmd.Flags().BoolVarP(&follow, "follow", "f", false, "wait for job to start, then stream logs")
}
