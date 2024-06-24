// Package k8s provides abstractions of some parts of the k8s API.
//
// The primary abstractions are built to simplify interactions with the batch API,
// specifically around managing jobs.
package k8s

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client exposes high-level Kubernetes API abstractions.
type Client interface {
	CreateJob(job *batchv1.Job) error
	DeleteJob(name string) error
	GetJob(name string) (*batchv1.Job, error)
	GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error)
	ListJobs() ([]batchv1.Job, error)
	GetJobEvents(name string) (string, error)
	GetPodEvents(name string) (string, error)
	GetPodsFromJob(jobName string) ([]string, error)
	GetJobFromPod(podName string) (string, error)
}

// NamespaceClient represents a namespaced Kubernetes API client.
type NamespaceClient struct {
	Clientset kubernetes.Interface
	Namespace string
}

// NewClient returns a Client the specified context and namespace.
func NewClient(context, namespace string) (Client, error) {
	config, namespace, err := buildClientConfig(context, namespace)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client := &NamespaceClient{clientset, namespace}

	return client, nil
}

// buildClientConfig returns a complete client config and the namespace for the given context.
func buildClientConfig(context, namespace string) (*rest.Config, string, error) {
	clientConfig := buildClientCmd(context, namespace)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not get k8s config for context %q: %w", context, err)
	}

	namespace, _, err = clientConfig.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("could not get namespace for context %q: %w", context, err)
	}

	return config, namespace, nil
}

// buildClientCmd returns an (incomplete) API server client config.
func buildClientCmd(context, namespace string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	if context != "" {
		overrides.CurrentContext = context
	}
	if namespace != "" {
		overrides.Context.Namespace = namespace
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
