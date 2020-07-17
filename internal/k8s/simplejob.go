package k8s

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleJob represents an extremely simplified k8s job specification.
type SimpleJob struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	WorkingDir string   `json:"workingDir,omitempty"`
	Command    []string `json:"command,omitempty"`

	Memory resource.Quantity `json:"memory,omitempty"`
	CPU    resource.Quantity `json:"cpu,omitempty"`
	GPU    resource.Quantity `json:"gpu,omitempty"`
}

var defaultVolumes = []corev1.Volume{{
	Name: "storage",
	VolumeSource: corev1.VolumeSource{
		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			ClaimName: "storage",
		},
	},
}}

var defaultVolumeMounts = []corev1.VolumeMount{{
	Name:      "storage",
	MountPath: "/storage",
}}

func (simple *SimpleJob) volumes() []corev1.Volume {
	return defaultVolumes
}

func (simple *SimpleJob) volumeMounts() []corev1.VolumeMount {
	return defaultVolumeMounts
}

func (simple *SimpleJob) resources() corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{Limits: corev1.ResourceList{
		"memory":         simple.Memory,
		"cpu":            simple.CPU,
		"nvidia.com/gpu": simple.GPU,
	}}

	return resources
}

func (simple *SimpleJob) containers() []corev1.Container {
	containers := []corev1.Container{{
		Name:         simple.Name,
		Image:        simple.Image,
		Command:      simple.Command,
		WorkingDir:   simple.WorkingDir,
		VolumeMounts: simple.volumeMounts(),
		Resources:    simple.resources(),

		Stdin: true,
		TTY:   true,
	}}

	return containers
}

func (simple *SimpleJob) meta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: simple.Name,
	}
}

// Expand expands the simplified job into a full job object.
func (simple *SimpleJob) Expand() *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: simple.meta(),
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{
				Containers: simple.containers(),
				Volumes:    simple.volumes(),
			}},
		},
	}

	return job
}
