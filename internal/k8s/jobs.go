package k8s

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

// SimpleJobSpec represents an extremely simplified k8s job specification.
type SimpleJobSpec struct {
	Name    string      `yaml:"name"`
	Image   string      `yaml:"image"`
	Command StringArray `yaml:"command"`
	GPU     uint        `yaml:"gpu"`
}

// DefaultLogOptions is the default set of options used for retrieving pod logs.
var DefaultLogOptions = &corev1.PodLogOptions{
	// TODO: Make these configurable via flags?
	Follow: true,
	// TailLines: int64Ptr(20),
}

// Expand expands the simplified job spec into a full job object.
func (spec *SimpleJobSpec) Expand() *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.Name,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: int32Ptr(0),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:    spec.Name,
							Image:   spec.Image,
							Command: spec.Command,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "storage",
								},
							},
						},
					},
				},
			},
		},
	}
}

// ListJobs returns all jobs.
func (kubectx *KubeContext) ListJobs() (*batchv1.JobList, error) {
	jobs, err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// DeleteJob deletes the job with the given name.
func (kubectx *KubeContext) DeleteJob(name string) error {
	deletePolicy := metav1.DeletePropagationBackground
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: int64Ptr(0),
		PropagationPolicy:  &deletePolicy,
	}

	err := kubectx.Client.BatchV1().Jobs(kubectx.Namespace).Delete(name, deleteOptions)
	return err
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
		return nil, fmt.Errorf("unable to get job: %v", err)
	}

	selector := labels.Set(job.Spec.Selector.MatchLabels).String()
	listOptions := metav1.ListOptions{LabelSelector: selector}
	pods, err := kubectx.Client.CoreV1().Pods(kubectx.Namespace).List(listOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get pods for job: %v", err)
	}

	if len(pods.Items) == 0 {
		return nil, nil
	}

	// TODO: Add support for multiple pods?
	pod := pods.Items[0]
	req := kubectx.Client.CoreV1().Pods(kubectx.Namespace).GetLogs(pod.Name, opts)
	return req, nil
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
