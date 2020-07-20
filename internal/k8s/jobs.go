package k8s

import (
	"fmt"
	"io/ioutil"
	"regexp"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

var defaultTerminationMessagePolicy = corev1.TerminationMessageFallbackToLogsOnError

var (
	// Do not restart failing jobs.
	defaultRestartPolicy = corev1.RestartPolicyOnFailure
	defaultBackoffLimit  = int32Ptr(0)
)

// DefaultLogOptions is the default set of options used when retrieving logs.
var DefaultLogOptions = &corev1.PodLogOptions{
	// TODO: Make these configurable via flags?
	Follow: true,
	// TailLines: int64Ptr(20),
}

// ListJobs returns all jobs.
func (kubectx *KubeContext) ListJobs() (*batchv1.JobList, error) {
	jobs, err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetJob returns the job with the given name.
func (kubectx *KubeContext) GetJob(name string) (*batchv1.Job, error) {
	getOptions := metav1.GetOptions{}
	job, err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).Get(name, getOptions)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return job, nil
}

// DeleteJob deletes the job with the given name.
func (kubectx *KubeContext) DeleteJob(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: int64Ptr(0),
		PropagationPolicy:  &deletePolicy,
	}

	err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).Delete(name, deleteOptions)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateJob creates a job with the given specification.
func (kubectx *KubeContext) CreateJob(job *batchv1.Job) error {
	_, err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).Create(job)
	return err
}

// GetJobLogs returns the pod logs for the job with the given name.
func (kubectx *KubeContext) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	getOptions := metav1.GetOptions{}
	job, err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).Get(name, getOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get job: %w", err)
	}

	selector := labels.Set(job.Spec.Selector.MatchLabels).String()
	listOptions := metav1.ListOptions{LabelSelector: selector}
	pods, err := kubectx.Client.CoreV1().Pods(kubectx.Namespace).List(listOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get pods for job: %w", err)
	}

	if len(pods.Items) == 0 {
		// TODO: Treat this as an error scenario?
		return nil, nil
	}

	// TODO: Add support for multiple pods?
	pod := pods.Items[0]
	req := kubectx.Client.CoreV1().Pods(kubectx.Namespace).GetLogs(pod.Name, opts)

	return req, nil
}

// ParseJob parses the file and returns the corresponding job.
func ParseJob(file string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// TODO: Refactor this by extracting functions, etc.
	var job *batchv1.Job
	re := regexp.MustCompile(`apiVersion:`)
	if re.Match(data) {
		job = &batchv1.Job{}
		if err := yaml.UnmarshalStrict(data, job); err != nil {
			return nil, err
		}
	} else {
		simple := &SimpleJob{}
		if err := yaml.UnmarshalStrict(data, simple); err != nil {
			return nil, err
		}
		job = simple.Expand()
	}

	return job, nil
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

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
