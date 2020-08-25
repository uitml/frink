// Package fake contains utilities intended for mocking out functionality in test code.
// Most notably, it contains mocks for interacting with the Kubernetes API.
package fake

import (
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// Client is a fake k8s.Client for interacting with a fake Kubernetes API, primarily intended for unit testing.
type Client struct {
	mock.Mock
}

// ListJobs returns a list of jobs based on the Jobs field in KubeClient.
func (client *Client) ListJobs() ([]batchv1.Job, error) {
	args := client.Called()
	jobs, _ := args.Get(0).([]batchv1.Job)

	return jobs, args.Error(1)
}

// CreateJob simulates creating a job.
func (client *Client) CreateJob(job *batchv1.Job) error {
	args := client.Called(job)

	return args.Error(0)
}

// DeleteJob simulates deleting a job.
func (client *Client) DeleteJob(name string) error {
	args := client.Called(name)

	return args.Error(0)
}

// GetJob searches through the items in the Jobs field in KubeClient, returning the first item with a matching name.
func (client *Client) GetJob(name string) (*batchv1.Job, error) {
	args := client.Called(name)
	job, _ := args.Get(0).(*batchv1.Job)

	return job, args.Error(1)
}

// GetJobLogs simulates returning a rest.Request that will stream logs for a the job with the matching name.
func (client *Client) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	args := client.Called(name, opts)
	req, _ := args.Get(0).(*rest.Request)

	return req, args.Error(1)
}
