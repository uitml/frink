package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootWithNoCommand(t *testing.T) {
	var buf strings.Builder
	rootCmd.SetOut(&buf)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.True(t, strings.Contains(buf.String(), "Usage:"))
}

func TestRootWithIncorrectCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"foo"})
	err := rootCmd.Execute()

	assert.Error(t, err)
}

func TestRootPersistentPreRun(t *testing.T) {
	before := client
	rootCmd.SetArgs([]string{"help"})
	err := rootCmd.Execute()
	after := client

	assert.NoError(t, err)
	assert.Nil(t, before, "expected package-scope 'client' variable to be uninitialized before command execution")
	assert.NotNil(t, after, "expected package-scope 'client' variable to have been initalized during command execution")
}
