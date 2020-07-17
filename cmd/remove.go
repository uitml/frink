package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove job from cluster",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job name must be specified")
		}

		name := args[0]
		if err := kubectx.DeleteJob(name); err != nil {
			return fmt.Errorf("unable to delete job: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
