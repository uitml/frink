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

func TestList_WithNoJobs(t *testing.T) {
	var out strings.Builder
	cmd := newListCmd()
	cmd.SetOut(&out)

	ctx := &listContext{
		CommandContext: cli.CommandContext{
			Client: &mocks.KubeClient{
				Jobs: []batchv1.Job{
					{
						ObjectMeta: v1.ObjectMeta{Name: "foo"},
					},
				},
			},
		},
	}

	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "foo")
}
