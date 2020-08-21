package k8s

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uitml/frink/internal/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
)

func TestListJobs(t *testing.T) {
	foo := newJob("foo")
	bar := newJob("bar")
	clientset := fake.NewSimpleClientset(&foo, &bar)
	client := NamespaceClient{
		Clientset: clientset,
	}

	jobs, err := client.ListJobs()
	assert.NoError(t, err)
	assert.Len(t, jobs.Items, 2)

	// Respond to job listing with an error.
	clientset.Fake.PrependReactor("list", "*", func(action kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("baz")
	})

	jobs, err = client.ListJobs()
	assert.Nil(t, jobs)
	assert.EqualError(t, err, "baz")
}

func TestGetJob(t *testing.T) {
	foo := newJob("foo")
	bar := newJob("bar")
	clientset := fake.NewSimpleClientset(&foo, &bar)
	client := NamespaceClient{
		Clientset: clientset,
	}

	job, err := client.GetJob("foo")
	assert.NoError(t, err)
	assert.Equal(t, &foo, job)

	job, err = client.GetJob("qux")
	assert.Nil(t, job)
	assert.NoError(t, err)

	// Respond to job retriebal with an error.
	clientset.Fake.PrependReactor("get", "*", func(action kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("baz")
	})

	job, err = client.GetJob("foo")
	assert.Nil(t, job)
	assert.EqualError(t, err, "baz")
}

func TestCreateJob(t *testing.T) {
	foo := newJob("foo")
	bar := newJob("bar")
	clientset := fake.NewSimpleClientset(&foo)
	client := NamespaceClient{
		Clientset: clientset,
	}

	err := client.CreateJob(&bar)
	assert.NoError(t, err)

	job, err := client.Clientset.BatchV1().Jobs("").Get("bar", v1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, &bar, job)
}

func TestDeleteJob(t *testing.T) {
	foo := newJob("foo")
	bar := newJob("bar")
	clientset := fake.NewSimpleClientset(&foo, &bar)
	client := NamespaceClient{
		Clientset: clientset,
	}

	before, _ := client.Clientset.BatchV1().Jobs("").List(v1.ListOptions{})
	err := client.DeleteJob("foo")
	after, _ := client.Clientset.BatchV1().Jobs("").List(v1.ListOptions{})

	assert.NoError(t, err)
	assert.Len(t, before.Items, 2)
	assert.Len(t, after.Items, 1)

	err = client.DeleteJob("foo")
	assert.NoError(t, err)

	// Respond to job deletion with an error.
	clientset.Fake.PrependReactor("delete", "*", func(action kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("baz")
	})

	err = client.DeleteJob("bar")
	assert.EqualError(t, err, "baz")
}

func TestOverrideJobSpec(t *testing.T) {
	job := newJob("foo", newZeroMemoryContainer())

	container := &job.Spec.Template.Spec.Containers[0]

	assert.Len(t, container.Resources.Limits, 2)
	assert.NotEqual(t, defaultTerminationMessagePolicy, container.TerminationMessagePolicy)
	assert.NotEqual(t, defaultBackoffLimit, job.Spec.BackoffLimit)
	assert.NotEqual(t, defaultRestartPolicy, job.Spec.Template.Spec.RestartPolicy)

	OverrideJobSpec(&job)

	assert.Len(t, container.Resources.Limits, 1)
	assert.Equal(t, defaultTerminationMessagePolicy, container.TerminationMessagePolicy)
	assert.Equal(t, defaultBackoffLimit, job.Spec.BackoffLimit)
	assert.Equal(t, defaultRestartPolicy, job.Spec.Template.Spec.RestartPolicy)
}

func TestRemoveZeroResources(t *testing.T) {
	container := newZeroMemoryContainer()

	qty, ok := container.Resources.Limits["memory"]
	assert.True(t, ok)
	assert.True(t, qty.IsZero())
	assert.Len(t, container.Resources.Limits, 2)

	removeZeroResources(&container)

	qty, ok = container.Resources.Limits["memory"]
	assert.False(t, ok)
	assert.True(t, qty.IsZero())
	assert.Len(t, container.Resources.Limits, 1)
}

func newJob(name string, containers ...corev1.Container) batchv1.Job {
	return batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Spec: batchv1.JobSpec{
			BackoffLimit: util.Int32Ptr(42),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers:    containers,
				},
			},
		},
	}
}

func newZeroMemoryContainer() corev1.Container {
	return corev1.Container{
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"cpu":    resource.MustParse("1"),
				"memory": resource.MustParse("0"),
			},
		},
	}
}
