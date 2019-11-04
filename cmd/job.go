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

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
