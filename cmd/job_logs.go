package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/pkg/kube/client"
	"github.com/uitml/frink/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var logsCmd = &cobra.Command{
	Use:     "logs [name]",
	Short:   "Fetch the logs of a job",
	Aliases: []string{"watch"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("job name must be specified")
		}

		name := args[0]

		clientset, namespace, err := client.ForContext("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %v", err)
		}

		getOptions := metav1.GetOptions{}
		job, err := clientset.BatchV1().Jobs(namespace).Get(name, getOptions)
		if err != nil {
			return fmt.Errorf("unable to get job: %v", err)
		}

		listOptions := metav1.ListOptions{
			LabelSelector: labels.Set(job.Spec.Selector.MatchLabels).String(),
		}
		pods, err := clientset.CoreV1().Pods(namespace).List(listOptions)
		if err != nil {
			return fmt.Errorf("unable to get pods for job: %v", err)
		}

		if len(pods.Items) == 0 {
			return nil
		}

		// TODO: Add support for multiple pods?
		pod := pods.Items[0]

		logOptions := &corev1.PodLogOptions{
			// TODO: Make these configurable via flags?
			TailLines: types.Int64Ptr(20),
			Follow:    true,
			// TODO: Figure out how to support both current and previous somehow (if necessary).
			// Previous: true,
		}
		req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, logOptions)
		if err != nil {
			return fmt.Errorf("unable to get logs: %v", err)
		}

		logs, err := req.Stream()
		if err != nil {
			return fmt.Errorf("unable to stream logs: %v", err)
		}
		defer logs.Close()

		r := bufio.NewReader(logs)
		for {
			p, err := r.ReadBytes('\n')
			if _, err := os.Stdout.Write(p); err != nil {
				return fmt.Errorf("unable to write output: %v", err)
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return fmt.Errorf("unable to read stream: %v", err)
			}
		}

		return nil
	},
}
