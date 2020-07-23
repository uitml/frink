// Package mocks contains utilities intended for mocking out functionality in test code.
// Most notably, it contains mocks for interacting with the Kubernetes API.
package mocks

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// KubeClient is a mock k8s.KubeClient for interacting with a fake Kubernetes API, primarily intended for unit testing.
type KubeClient struct {
	Jobs []batchv1.Job
}

// ListJobs returns a list of jobs based on the Jobs field in KubeClient.
func (client *KubeClient) ListJobs() (*batchv1.JobList, error) {
	jobs := &batchv1.JobList{
		Items: client.Jobs,
	}

	return jobs, nil
}

// CreateJob simulates creating a job.
func (client *KubeClient) CreateJob(job *batchv1.Job) error {
	return nil
}

// DeleteJob simulates deleting a job.
func (client *KubeClient) DeleteJob(name string) error {
	return nil
}

// GetJob searches through the items in the Jobs field in KubeClient, returning the first item with a matching name.
func (client *KubeClient) GetJob(name string) (*batchv1.Job, error) {
	return nil, nil
}

// GetJobLogs simulates returning a rest.Request that will stream logs for a the job with the matching name.
func (client *KubeClient) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	return nil, nil
}
