package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
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

		kubectx, err := k8s.Client("")
		if err != nil {
			return fmt.Errorf("unable to get kube client: %v", err)
		}

		req, err := kubectx.GetJobLogs(name, k8s.DefaultLogOptions)
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
