package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s"
	"github.com/uitml/frink/internal/k8s/fake"
)

// Top-level functionality.

func TestRunPreRun(t *testing.T) {
	ctx := &runContext{}
	cmd := newRunCmd()

	assert.Nil(t, ctx.Client)
	ctx.PreRun(cmd, []string{})
	assert.NotNil(t, ctx.Client)
}

func TestRunRunNewJob(t *testing.T) {
	var out strings.Builder
	cmd := newRunCmd()
	cmd.SetOut(&out)

	client := &fake.Client{}

	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := k8s.NewJobParser(fs)

	ctx := &runContext{
		CommandContext: cli.CommandContext{
			Out:    cmd.OutOrStderr(),
			Err:    cmd.ErrOrStderr(),
			Client: client,
		},
		JobParser: parser,
	}

	filename := "job.yaml"
	job, _ := parser.Parse(filename)
	k8s.OverrideJobSpec(job)

	client.On("GetJob", job.Name).Return(nil, nil)
	client.On("CreateJob", job).Return(nil)

	err := ctx.Run(cmd, []string{filename})
	assert.NoError(t, err)
}

func TestRunRunExistingJob(t *testing.T) {
	var out strings.Builder
	cmd := newRunCmd()
	cmd.SetOut(&out)

	client := &fake.Client{}

	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := k8s.NewJobParser(fs)

	ctx := &runContext{
		CommandContext: cli.CommandContext{
			Out:    cmd.OutOrStderr(),
			Err:    cmd.ErrOrStderr(),
			Client: client,
		},
		JobParser: parser,
	}

	filename := "job.yaml"
	job, _ := parser.Parse(filename)
	k8s.OverrideJobSpec(job)

	client.On("GetJob", job.Name).Return(job, nil).Once() // Only return once to emulate deletion
	client.On("GetJob", job.Name).Return(nil, nil)
	client.On("DeleteJob", job.Name).Return(nil)
	client.On("CreateJob", job).Return(nil)

	err := ctx.Run(cmd, []string{filename})
	assert.NoError(t, err)

	client.AssertExpectations(t)
}

func TestRunRunMissingArgument(t *testing.T) {
	var out strings.Builder
	cmd := newRunCmd()
	cmd.SetOut(&out)

	client := &fake.Client{}

	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := k8s.NewJobParser(fs)

	ctx := &runContext{
		CommandContext: cli.CommandContext{
			Out:    cmd.OutOrStderr(),
			Err:    cmd.ErrOrStderr(),
			Client: client,
		},
		JobParser: parser,
	}

	err := ctx.Run(cmd, []string{})
	assert.EqualError(t, err, "job specification file must be specified")

	client.AssertExpectations(t)
}

func TestRunRunCreationError(t *testing.T) {
	var out strings.Builder
	cmd := newRunCmd()
	cmd.SetOut(&out)

	client := &fake.Client{}

	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := k8s.NewJobParser(fs)

	ctx := &runContext{
		CommandContext: cli.CommandContext{
			Out:    cmd.OutOrStderr(),
			Err:    cmd.ErrOrStderr(),
			Client: client,
		},
		JobParser: parser,
	}

	filename := "job.yaml"
	job, _ := parser.Parse(filename)
	k8s.OverrideJobSpec(job)

	client.On("GetJob", job.Name).Return(nil, nil)
	client.On("CreateJob", job).Return(errors.New("baz"))

	err := ctx.Run(cmd, []string{filename})
	assert.EqualError(t, err, "unable to create job: baz")

	client.AssertExpectations(t)
}
