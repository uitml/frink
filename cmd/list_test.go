package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/mocks"
)

func TestListWithNoJobs(t *testing.T) {
	var buf strings.Builder
	cmd := NewListCmd()
	cmd.SetOut(&buf)

	client := &mocks.KubeClient{}
	ctx := &ListContext{Client: client}
	err := ctx.Run(cmd, []string{})

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "foo")
}
