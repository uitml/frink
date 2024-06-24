package k8s

import (
	"context"
	"fmt"
	"strings"

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
	Follow: true,
}

// ListJobs returns all jobs.
func (client *NamespaceClient) ListJobs() ([]batchv1.Job, error) {
	jobs, err := client.Clientset.BatchV1().Jobs(client.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs.Items, nil
}

// GetJob returns the job with the given name.
func (client *NamespaceClient) GetJob(name string) (*batchv1.Job, error) {
	getOptions := metav1.GetOptions{}
	job, err := client.Clientset.BatchV1().Jobs(client.Namespace).Get(context.TODO(), name, getOptions)
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
	deleteOptions := metav1.DeleteOptions{
		GracePeriodSeconds: util.Int64Ptr(0),
		PropagationPolicy:  &deletePolicy,
	}

	err := client.Clientset.BatchV1().Jobs(client.Namespace).Delete(context.TODO(), name, deleteOptions)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateJob creates a job with the given specification.
func (client *NamespaceClient) CreateJob(job *batchv1.Job) error {
	opts := metav1.CreateOptions{}
	_, err := client.Clientset.BatchV1().Jobs(client.Namespace).Create(context.TODO(), job, opts)
	return err
}

// GetJobLogs returns the pod logs for the job with the given name.
func (client *NamespaceClient) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	getOptions := metav1.GetOptions{}
	job, err := client.Clientset.BatchV1().Jobs(client.Namespace).Get(context.TODO(), name, getOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get job: %w", err)
	}

	selector := labels.Set(job.Spec.Selector.MatchLabels).String()
	listOptions := metav1.ListOptions{LabelSelector: selector}
	pods, err := client.Clientset.CoreV1().Pods(client.Namespace).List(context.TODO(), listOptions)
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

func (client *NamespaceClient) GetPodEvents(podName string) (string, error) {
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Pod", podName),
	}
	events, err := client.Clientset.CoreV1().Events(client.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return "", fmt.Errorf("unable to list events for pod %s: %w", podName, err)
	}

	var eventDetails strings.Builder
	for _, event := range events.Items {
		eventDetails.WriteString(fmt.Sprintf("%s %s: %s\n", event.CreationTimestamp, event.Reason, event.Message))
	}

	return eventDetails.String(), nil
}

func (client *NamespaceClient) GetJobEvents(jobName string) (string, error) {
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Job", jobName),
	}
	events, err := client.Clientset.CoreV1().Events(client.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return "", fmt.Errorf("unable to list events for job %s: %w", jobName, err)
	}

	var eventDetails strings.Builder
	for _, event := range events.Items {
		eventDetails.WriteString(fmt.Sprintf("%s %s: %s\n", event.CreationTimestamp, event.Reason, event.Message))
	}

	return eventDetails.String(), nil
}

func (client *NamespaceClient) GetPodsFromJob(jobName string) ([]string, error) {
	job, err := client.GetJob(jobName)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("job %s not found", jobName)
	}

	selector := labels.Set(job.Spec.Selector.MatchLabels).String()
	listOptions := metav1.ListOptions{LabelSelector: selector}
	pods, err := client.Clientset.CoreV1().Pods(client.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

func (client *NamespaceClient) GetJobFromPod(podName string) (string, error) {
	pod, err := client.Clientset.CoreV1().Pods(client.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if pod == nil {
		return "", fmt.Errorf("pod %s not found", podName)
	}

	// Assuming the job name is stored in the pod's labels under a specific key
	jobName, ok := pod.Labels["job-name"]
	if !ok {
		return "", fmt.Errorf("no job associated with pod %s", podName)
	}

	return jobName, nil
}
