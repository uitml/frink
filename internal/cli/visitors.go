// Package cli contains various shared command-line functionality, such as command algorithms, command context, and configuration.
package cli

import "github.com/spf13/cobra"

// VisitAll visits the entire command tree rooted at cmd, invoking fn on each command visited.
func VisitAll(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, child := range cmd.Commands() {
		VisitAll(child, fn)
	}
}

// DisableFlagsInUseLine enables the DisableFlagsInUseLine flag on the entire command tree, rooted at cmd.
func DisableFlagsInUseLine(cmd *cobra.Command) {
	VisitAll(cmd, func(c *cobra.Command) {
		c.DisableFlagsInUseLine = true
	})
}
