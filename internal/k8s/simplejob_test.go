package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandSetsObjectMeta(t *testing.T) {
	simple := &SimpleJob{Name: "foo"}

	job := simple.Expand()
	assert.Equal(t, simple.Name, job.Name)
}

func TestExpandDefinesContainer(t *testing.T) {
	simple := &SimpleJob{
		Name:       "foo",
		Image:      "ubuntu:latest",
		Command:    []string{"echo", "hello world"},
		WorkingDir: "/storage",
	}

	job := simple.Expand()
	pod := job.Spec.Template.Spec

	container := pod.Containers[0]
	assert.Equal(t, simple.Name, container.Name)
	assert.Equal(t, simple.Image, container.Image)
	assert.Equal(t, simple.Command, container.Command)
	assert.Equal(t, simple.WorkingDir, container.WorkingDir)

	assert.True(t, container.Stdin)
	assert.True(t, container.TTY)
}

func TestExpandDefinesVolumes(t *testing.T) {
	simple := &SimpleJob{}
	job := simple.Expand()

	pod := job.Spec.Template.Spec
	assert.Equal(t, defaultVolumes, pod.Volumes)

	container := pod.Containers[0]
	assert.Equal(t, defaultVolumeMounts, container.VolumeMounts)
}

func TestDefaultVolumeCompatibility(t *testing.T) {
	var claimNames []string
	for _, volume := range defaultVolumes {
		claimNames = append(claimNames, volume.VolumeSource.PersistentVolumeClaim.ClaimName)
	}

	var mountNames []string
	for _, mount := range defaultVolumeMounts {
		mountNames = append(mountNames, mount.Name)
	}

	assert.Equal(t, claimNames, mountNames)
}
