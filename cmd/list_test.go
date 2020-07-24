package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/mocks"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	activeJob = batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: batchv1.JobStatus{
			Active: 1,
		},
	}

	failedJob = batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: batchv1.JobStatus{
			Failed: 1,
		},
	}

	stoppedJob = batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       batchv1.JobSpec{Completions: nil},
		Status: batchv1.JobStatus{
			Succeeded: 0,
		},
	}

	successfulJob = batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: batchv1.JobStatus{
			Succeeded: 1,
		},
	}
)

func TestListPreRun(t *testing.T) {
	ctx := &listContext{}
	cmd := newListCmd()
	assert.Nil(t, ctx.Client)

	ctx.PreRun(cmd, []string{})
	assert.NotNil(t, ctx.Client)
}

func TestListRunBrokenClient(t *testing.T) {
	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: &mocks.KubeClient{
				Err: errors.New("foo"),
			},
		},
	}

	cmd := newListCmd()
	err := ctx.Run(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foo")
}

func TestListOutputWithNoJobs(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: &mocks.KubeClient{},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Equal(t, 1, strings.Count(out.String(), "\n"))
}

func TestListOutputWithJobs(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: &mocks.KubeClient{
				Jobs: []batchv1.Job{
					successfulJob,
				},
			},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
	assert.Contains(t, out.String(), "Succeeded")
}

func TestStatusActiveJob(t *testing.T) {
	job := activeJob
	out := status(job)

	assert.Equal(t, out, "Active")
}

func TestStatusFailedJob(t *testing.T) {
	job := failedJob
	out := status(job)

	assert.Equal(t, out, "Failed")
}

func TestStatusStoppedJob(t *testing.T) {
	job := stoppedJob
	out := status(job)

	assert.Equal(t, out, "Stopped")
}

func TestStatusSuccessfulJob(t *testing.T) {
	job := successfulJob
	out := status(job)

	assert.Equal(t, "Succeeded", out)
}

func TestHeaderTrailingTab(t *testing.T) {
	out := header()

	assert.Regexp(t, "\t$", out)
}

func TestRowTrailingTab(t *testing.T) {
	job := successfulJob
	out := row(job)

	assert.Regexp(t, "\t$", out)
}

func TestMatchingTabCount(t *testing.T) {
	job := successfulJob
	rowOut := row(job)
	hdrOut := header()

	assert.Equal(t, strings.Count(rowOut, "\t"), strings.Count(hdrOut, "\t"))
}
