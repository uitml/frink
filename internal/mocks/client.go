package mocks

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type KubeClient struct {
	Jobs []batchv1.Job
}

func (client *KubeClient) ListJobs() (*batchv1.JobList, error) {
	jobs := &batchv1.JobList{
		Items: client.Jobs,
	}

	return jobs, nil
}

func (client *KubeClient) CreateJob(job *batchv1.Job) error {
	return nil
}

func (client *KubeClient) DeleteJob(name string) error {
	return nil
}

func (client *KubeClient) GetJob(name string) (*batchv1.Job, error) {
	return nil, nil
}

func (client *KubeClient) GetJobLogs(name string, opts *corev1.PodLogOptions) (*rest.Request, error) {
	return nil, nil
}
