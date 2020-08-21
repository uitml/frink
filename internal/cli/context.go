package cli

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/k8s"
)

// CommandContext is a common context facility used by all commands in the cmd package that interact with the Kubernetes API.
type CommandContext struct {
	CommandInitializer

	// Out can be used to write to an output stream, typically stdout.
	Out io.Writer

	// Err can be used to write to an error output stream, typically stderr.
	Err io.Writer

	// Client can be used for interacting with the Kubernetes API.
	Client k8s.Client
}

// CommandInitializer is an interface that is used to initialize a CommandContext.
//
// The most common thing to do during initialization is to set the Client (k8s.KubeClient) field.
type CommandInitializer interface {
	// Initialize initializes the CommandContext, optionally using the cobra.Command flags, etc.
	Initialize(cmd *cobra.Command) error
}

// Initialize is used to initialize the Client field on the CommandContext.
//
// The initialization is performed by first getting user configuration, where some settings might be overriden by command-line flags.
// Then a k8s.KubeClient is created using the context and namespace specified by the user configuration.
// Finally, the Client field on the CommandContext is set to the newly created k8s.KubeClient.
func (ctx *CommandContext) Initialize(cmd *cobra.Command) error {
	cfg, err := ParseConfig(cmd)
	if err != nil {
		return err
	}

	client, err := k8s.NewClient(cfg.Context, cfg.Namespace)
	if err != nil {
		return err
	}

	ctx.Out = cmd.OutOrStderr()
	ctx.Err = cmd.ErrOrStderr()
	ctx.Client = client

	return nil
}
