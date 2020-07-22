package cli

import (
	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
)

type CommandContext struct {
	CommandInitializer
	Client k8s.KubeClient
}

type CommandInitializer interface {
	Initialize(*cobra.Command) error
}

func (ctx *CommandContext) Initialize(cmd *cobra.Command) error {
	cfg, err := ParseConfig(cmd)
	if err != nil {
		return err
	}

	client, err := k8s.Client(cfg.Context, cfg.Namespace)
	if err != nil {
		return err
	}

	ctx.Client = client

	return nil
}
