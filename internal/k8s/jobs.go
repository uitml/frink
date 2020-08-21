package k8s

import (
	"fmt"

	"github.com/uitml/frink/internal/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

var defaultTerminationMessagePolicy = corev1.TerminationMessageFallbackToLogsOnError

var (
	// Do not restart failing jobs.
	defaultRestartPolicy = corev1.RestartPolicyOnFailure
	defaultBackoffLimit  = util.Int32Ptr(0)
)

// DefaultLogOptions is the default set of options used when retrieving logs.
var DefaultLogOptions = &corev1.PodLogOptions{
	// TODO: Make these configurable via flags?
	Follow: true,
	// TailLines: int64Ptr(20),
}

// ListJobs returns all jobs.
func (client *NamespaceClient) ListJobs() ([]batchv1.Job, error) {
	jobs, err := client.Clientset.BatchV1().Jobs(client.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs.Items, nil
}

// GetJob returns the job with the given name.
func (client *NamespaceClient) GetJob(name string) (*batchv1.Job, error) {
	getOptions := metav1.GetOptions{}
	job, err := client.Clientset.BatchV1().Jobs(client.Namespace).Get(name, getOptions)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return job, nil
}

// DeleteJob deletes the job with the given name.
func (client *NamespaceClient) DeleteJob(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: util.Int64Ptr(0),
		PropagationPolicy:  &deletePolicy,
	}

	err := client.Clientset.BatchV1().Jobs(client.Namespace).Delete(name, deleteOptions)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateJob creates a job with the given specification.
func (client *NamespaceClient) CreateJob(job *batchv1.Job) error {
	_, err := client.Clientset.BatchV1().Jobs(client.Namespace).Create(job)
	return err
}

// GetJobLogs returns the pod logs for the job with the given name.
func (client *NamespaceClient) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	getOptions := metav1.GetOptions{}
	job, err := client.Clientset.BatchV1().Jobs(client.Namespace).Get(name, getOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get job: %w", err)
	}

	selector := labels.Set(job.Spec.Selector.MatchLabels).String()
	listOptions := metav1.ListOptions{LabelSelector: selector}
	pods, err := client.Clientset.CoreV1().Pods(client.Namespace).List(listOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get pods for job: %w", err)
	}

	if len(pods.Items) == 0 {
		// TODO: Treat this as an error scenario?
		return nil, nil
	}

	// TODO: Add support for multiple pods?
	pod := pods.Items[0]
	req := client.Clientset.CoreV1().Pods(client.Namespace).GetLogs(pod.Name, opts)

	return req, nil
}

// OverrideJobSpec removes zero quantity resources, and sets other important defaults.
func OverrideJobSpec(job *batchv1.Job) {
	containers := job.Spec.Template.Spec.Containers
	for i := range containers {
		container := &containers[i]
		removeZeroResources(container)
		setTerminationPolicy(container)
	}

	setRestartPolicy(job)
}

func removeZeroResources(container *corev1.Container) {
	limits := container.Resources.Limits
	for k, v := range limits {
		if v.IsZero() {
			// TODO: Notify user.
			delete(limits, k)
		}
	}
}

func setTerminationPolicy(container *corev1.Container) {
	container.TerminationMessagePolicy = defaultTerminationMessagePolicy
}

func setRestartPolicy(job *batchv1.Job) {
	job.Spec.BackoffLimit = defaultBackoffLimit
	job.Spec.Template.Spec.RestartPolicy = defaultRestartPolicy
}
