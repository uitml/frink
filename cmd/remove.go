package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	waitForDelete bool
)

var removeCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove job from cluster",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job name must be specified")
		}

		name := args[0]

		job, err := client.GetJob(name)
		if err != nil {
			return fmt.Errorf("unable to get job: %w", err)
		}

		if job == nil {
			fmt.Printf("Nothing to delete: no job named %s found\n", name)
			return nil
		}

		fmt.Printf("Deleting job %s...\n", name)
		if err := client.DeleteJob(name); err != nil {
			return fmt.Errorf("unable to delete job: %w", err)
		}

		if !waitForDelete {
			return nil
		}

		if err := waitUntilDeleted(name); err != nil {
			return fmt.Errorf("timed out waiting for job to be deleted: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	flags := removeCmd.Flags()
	flags.BoolVarP(&waitForDelete, "wait", "w", false, "wait for job to be deleted")
}
