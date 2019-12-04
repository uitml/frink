package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
)

var removeCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove job from cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job name must be specified")
		}

		name := args[0]

		kubectx, err := k8s.Client("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %w", err)
		}

		err = kubectx.DeleteJob(name)
		if err != nil {
			return fmt.Errorf("unable to delete job: %w", err)
		}

		return nil
	},
}
