package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/mocks"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
					{
						ObjectMeta: v1.ObjectMeta{Name: "foo"},
						Status: batchv1.JobStatus{
							Succeeded: 1,
						},
					},
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
	job := batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: batchv1.JobStatus{
			Active: 1,
		},
	}

	out := status(job)

	assert.Equal(t, out, "Active")
}

func TestStatusFailedJob(t *testing.T) {
	job := batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: batchv1.JobStatus{
			Failed: 1,
		},
	}

	out := status(job)

	assert.Equal(t, out, "Failed")
}

func TestStatusStoppedJob(t *testing.T) {
	job := batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       batchv1.JobSpec{Completions: nil},
		Status: batchv1.JobStatus{
			Succeeded: 0,
		},
	}

	out := status(job)

	assert.Equal(t, out, "Stopped")
}

func TestStatusSuccessfulJob(t *testing.T) {
	job := batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       batchv1.JobSpec{Completions: nil},
		Status: batchv1.JobStatus{
			Succeeded: 1,
		},
	}

	out := status(job)

	assert.Equal(t, "Succeeded", out)
}
