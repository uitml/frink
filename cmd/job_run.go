package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/pkg/kube/client"
	"github.com/uitml/frink/pkg/kube/retry"
	"github.com/uitml/frink/pkg/types"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Schedule a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job specification file must be specified")
		}

		file := args[0]
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("specified file does not exist: %v", file)
			}

			return fmt.Errorf("unable to access file: %v", err)
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("unable to read file: %v", err)
		}

		// TODO: Implement support for "normal" k8s job specs. Determine based on presence of 'kind' and/or 'apiVersion' keys?
		spec := types.SimpleJobSpec{}
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return fmt.Errorf("unable to parse file: %v", err)
		}
		job := spec.Expand()

		clientset, namespace, err := client.ForContext("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %v", err)
		}

		jobClient := clientset.BatchV1().Jobs(namespace)

		// Delete existing job with same name.
		deletePolicy := metav1.DeletePropagationBackground
		deleteOptions := &metav1.DeleteOptions{
			GracePeriodSeconds: types.Int64Ptr(0),
			PropagationPolicy:  &deletePolicy,
		}
		err = jobClient.Delete(spec.Name, deleteOptions)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("unable to previous job: %v", err)
		}

		// Try to create the job using retry with backoff.
		err = retry.RetryOnExists(retry.DefaultBackoff, func() error {
			_, err = jobClient.Create(job)
			return err
		})
		if err != nil {
			return fmt.Errorf("unable to create job: %v", err)
		}

		// TODO: Implement support for streaming job/pod log to stdout.

		return nil
	},
}
