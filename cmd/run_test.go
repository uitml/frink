package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/cli"
	"github.com/uitml/frink/internal/k8s/fake"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Top-level functionality.

func TestRunPreRun(t *testing.T) {
	ctx := &runContext{}
	cmd := newRunCmd()

	assert.Nil(t, ctx.Client)
	ctx.PreRun(cmd, []string{})
	assert.NotNil(t, ctx.Client)
}

func TestRunRun(t *testing.T) {
	var out strings.Builder
	cmd := newRunCmd()
	cmd.SetOut(&out)

	client := &fake.Client{}
	parser := &fake.JobParser{}

	ctx := &runContext{
		CommandContext: cli.CommandContext{
			Out:    cmd.OutOrStderr(),
			Err:    cmd.ErrOrStderr(),
			Client: client,
		},
		JobParser: parser,
	}

	job := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
	}

	client.On("GetJob", job.Name).Return(nil, nil)
	client.On("CreateJob", job).Return(nil)
	parser.On("Parse").Return(job, nil)

	err := ctx.Run(cmd, []string{"bar"})
	assert.NoError(t, err)
}
