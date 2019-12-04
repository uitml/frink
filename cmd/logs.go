package cmd

import (
	"bufio"
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
			return fmt.Errorf("unable to get kube client: %w", err)
		}

		req, err := kubectx.GetJobLogs(name, k8s.DefaultLogOptions)
		if err != nil {
			return fmt.Errorf("unable to get logs: %w", err)
		}

		stream, err := req.Stream()
		if err != nil {
			return fmt.Errorf("unable to stream logs: %w", err)
		}
		defer stream.Close()

		reader := bufio.NewReader(stream)
		if _, err := io.Copy(os.Stdout, reader); err != nil {
			return fmt.Errorf("unable to write output: %w", err)
		}

		return nil
	},
}
