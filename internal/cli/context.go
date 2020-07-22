package cli

import (
	"github.com/spf13/viper"
	"github.com/uitml/frink/internal/k8s"
)

type CommandContext interface {
	SetClient(k8s.KubeClient)
}

func Initialize(ctx CommandContext) error {
	context := viper.GetString("context")
	namespace := viper.GetString("namespace")

	client, err := k8s.Client(context, namespace)
	if err != nil {
		return err
	}

	ctx.SetClient(client)

	return nil
}
