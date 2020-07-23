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

func TestList_NoJobs(t *testing.T) {
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

func TestList_OneSuccessfulJob(t *testing.T) {
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

func TestList_OneFailedJob(t *testing.T) {
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
							Failed: 1,
						},
					},
				},
			},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
	assert.Contains(t, out.String(), "Failed")
}

func TestList_OneActiveJob(t *testing.T) {
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
							Active: 1,
						},
					},
				},
			},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
	assert.Contains(t, out.String(), "Active")
}

func TestList_OneStoppedJob(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: &mocks.KubeClient{
				Jobs: []batchv1.Job{
					{
						ObjectMeta: v1.ObjectMeta{Name: "foo"},
						Spec:       batchv1.JobSpec{Completions: nil},
						Status: batchv1.JobStatus{
							Succeeded: 0,
						},
					},
				},
			},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
	assert.Contains(t, out.String(), "Stopped")
}
