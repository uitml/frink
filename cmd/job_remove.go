package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/pkg/kube/client"
	"github.com/uitml/frink/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var removeCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job name must be specified")
		}

		name := args[0]

		clientset, namespace, err := client.ForContext("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %v", err)
		}

		deletePolicy := metav1.DeletePropagationBackground
		deleteOptions := &metav1.DeleteOptions{
			GracePeriodSeconds: types.Int64Ptr(0),
			PropagationPolicy:  &deletePolicy,
		}
		err = clientset.BatchV1().Jobs(namespace).Delete(name, deleteOptions)
		if err != nil {
			return fmt.Errorf("unable to delete job: %v", err)
		}

		return nil
	},
}
