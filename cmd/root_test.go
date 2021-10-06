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
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"foo"})
	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, out.String(), "unknown command")
}

func TestRootFlagsInUseLineDisabled(t *testing.T) {
	cmd := newRootCmd()

	for _, c := range cmd.Commands() {
		assert.Truef(t, c.DisableFlagsInUseLine, "Command: %s", c.Name())
	}
}
