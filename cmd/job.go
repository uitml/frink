package cmd

import "github.com/spf13/cobra"

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage scheduled jobs",
}

func init() {
	jobCmd.AddCommand(listCmd)
	jobCmd.AddCommand(logsCmd)
	jobCmd.AddCommand(removeCmd)
	jobCmd.AddCommand(runCmd)
}
