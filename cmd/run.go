package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s"
	"github.com/uitml/frink/internal/k8s/retry"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

var backoff = wait.Backoff{
	Duration: 100 * time.Millisecond,
	Factor:   1.0,
	Steps:    1200,
}

type runContext struct {
	cli.CommandContext

	Follow bool
}

func newRunCmd() *cobra.Command {
	ctx := &runContext{}
	cmd := &cobra.Command{
		Use:   "run <file>",
		Short: "Schedule a job on the cluster",

		PreRunE: ctx.PreRun,
		RunE:    ctx.Run,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&ctx.Follow, "follow", "f", false, "wait for job to start, then stream logs")

	return cmd
}

func (ctx *runContext) PreRun(cmd *cobra.Command, args []string) error {
	return ctx.Initialize(cmd)
}

func (ctx *runContext) Run(cmd *cobra.Command, args []string) error {
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

	if err := ctx.DeletePreviousJob(job.Name); err != nil {
		return fmt.Errorf("unable to delete previous job: %w", err)
	}

	// Try to create the job using retry.
	// This handles scenarios where an existing job is still being terminated, etc.
	fmt.Println("Creating job...")
	err = retry.OnExists(backoff, func() error { return ctx.Client.CreateJob(job) })
	if err != nil {
		return fmt.Errorf("unable to create job: %w", err)
	}

	if !ctx.Follow {
		return nil
	}

	if err := ctx.WaitUntilJobStarted(job.Name); err != nil {
		return fmt.Errorf("timed out waiting for job to start: %w", err)
	}

	// TODO: Ensure nil references are properly handled in this block.
	err = retry.OnError(backoff, apierrors.IsBadRequest, func() error {
		req, err := ctx.Client.GetJobLogs(job.Name, k8s.DefaultLogOptions)
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
}

func (ctx *runContext) DeletePreviousJob(name string) error {
	job, err := ctx.Client.GetJob(name)
	if err != nil {
		return err
	}

	if job != nil {
		fmt.Println("Deleting previous job...")
		err = ctx.Client.DeleteJob(job.Name)
		if err != nil {
			return err
		}

		if err := ctx.WaitUntilJobDeleted(job.Name); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *runContext) WaitUntilJobDeleted(name string) error {
	err := wait.Poll(100*time.Millisecond, 120*time.Second, func() (bool, error) {
		job, err := ctx.Client.GetJob(name)
		if err != nil {
			return false, err
		}

		return job == nil, nil
	})

	return err
}

func (ctx *runContext) WaitUntilJobStarted(name string) error {
	err := wait.Poll(100*time.Millisecond, 120*time.Second, func() (bool, error) {
		job, err := ctx.Client.GetJob(name)
		if err != nil {
			return false, err
		}

		return job != nil && job.Status.Active > 0, nil
	})

	return err
}
