package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeContext represents a concrete Kubernetes API context.
type KubeContext struct {
	Client    *kubernetes.Clientset
	Namespace string
}

// Client returns a k8s client and namespace for the specified context.
func Client(context string) (*KubeContext, error) {
	config, namespace, err := buildClientConfig(context)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubectx := &KubeContext{client, namespace}

	return kubectx, nil
}

// buildClientConfig returns a complete client config and the namespace for the given context.
func buildClientConfig(context string) (*rest.Config, string, error) {
	clientConfig := buildClientCmd(context)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not get k8s config for context %q: %w", context, err)
	}

	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("could not get namespace for context %q: %w", context, err)
	}

	return config, namespace, nil
}

// buildClientCmd returns an (incomplete) API server client config.
func buildClientCmd(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
