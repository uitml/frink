package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Top-level functionality.

func TestRunPreRun(t *testing.T) {
	ctx := &runContext{}
	cmd := newRunCmd()

	assert.Nil(t, ctx.Client)
	ctx.PreRun(cmd, []string{})
	assert.NotNil(t, ctx.Client)
}
