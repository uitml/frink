package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootNoCommand(t *testing.T) {
	var out strings.Builder
	cmd := newRootCmd()
	cmd.SetOut(&out)
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "Usage:")
}

func TestRootUnknownCommand(t *testing.T) {
	var out strings.Builder
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"foo"})
	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, out.String(), "unknown command")
}
