package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/mock"
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

func jobList(jobs ...batchv1.Job) *batchv1.JobList {
	return &batchv1.JobList{Items: jobs}
}

// Top-level functionality.

func TestListPreRun(t *testing.T) {
	ctx := &listContext{}
	cmd := newListCmd()

	assert.Nil(t, ctx.Client)
	ctx.PreRun(cmd, []string{})
	assert.NotNil(t, ctx.Client)
}

func TestListRunBrokenClient(t *testing.T) {
	client := &mock.KubeClient{}
	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: client,
		},
	}

	client.On("ListJobs").Return(nil, errors.New("foo"))

	cmd := newListCmd()
	err := ctx.Run(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foo")

	client.AssertExpectations(t)
}

// Raw command output

func TestListOutputWithNoJobs(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	client := &mock.KubeClient{}

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: client,
		},
	}

	client.On("ListJobs").Return(jobList(), nil)

	err := ctx.Run(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, strings.Count(out.String(), "\n"))

	client.AssertExpectations(t)
}

func TestListOutputWithJobs(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	client := &mock.KubeClient{}

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: client,
		},
	}

	client.On("ListJobs").Return(jobList(successfulJob), nil)

	err := ctx.Run(cmd, []string{})
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
	assert.Contains(t, out.String(), "Succeeded")

	client.AssertExpectations(t)
}

// Table formatting

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

// Column-level formatting

func TestStatusActiveJob(t *testing.T) {
	job := activeJob
	out := status(job)
	assert.Equal(t, "Active", out)
}

func TestStatusFailedJob(t *testing.T) {
	job := failedJob
	out := status(job)
	assert.Equal(t, "Failed", out)
}

func TestStatusStoppedJob(t *testing.T) {
	job := stoppedJob
	out := status(job)
	assert.Equal(t, "Stopped", out)
}

func TestStatusSuccessfulJob(t *testing.T) {
	job := successfulJob
	out := status(job)
	assert.Equal(t, "Succeeded", out)
}

func TestCompletionsActiveJob(t *testing.T) {
	job := activeJob
	out := completions(job)
	assert.Equal(t, "0/1", out)
}

func TestCompletionsFailedJob(t *testing.T) {
	job := failedJob
	out := completions(job)
	assert.Equal(t, "0/1", out)
}

func TestCompletionsStoppedJob(t *testing.T) {
	job := stoppedJob
	out := completions(job)
	assert.Equal(t, "0/0", out)
}

func TestComplectionsSuccessfulJob(t *testing.T) {
	job := successfulJob
	out := completions(job)
	assert.Equal(t, "1/1", out)
}

func TestComplectionsMultiplePods(t *testing.T) {
	job := activeJob
	job.Status.Succeeded = 1

	out := completions(job)
	assert.Equal(t, "1/2", out)
}

func TestDuration(t *testing.T) {
	now := time.Now()
	d, _ := time.ParseDuration("1h2m3ms4ns")
	past := now.Add(-d)

	job := activeJob
	job.Status.StartTime = &v1.Time{Time: past}

	out := duration(job)
	assert.Equal(t, "1 hour 2 minutes", out)
}

func TestAgeActiveJob(t *testing.T) {
	now := time.Now()
	d, _ := time.ParseDuration("2m")
	past := now.Add(-d)

	job := activeJob
	job.Status.StartTime = &v1.Time{Time: past}

	out := age(job)
	assert.Equal(t, "2 minutes ago", out)
}

func TestAgeSuccessfulJob(t *testing.T) {
	now := time.Now()
	d, _ := time.ParseDuration("1h")
	past := now.Add(-d)

	job := successfulJob
	job.Status.StartTime = &v1.Time{Time: past}
	job.Status.CompletionTime = &v1.Time{Time: now}

	out := age(job)
	assert.Equal(t, "1 hour ago", out)
}
