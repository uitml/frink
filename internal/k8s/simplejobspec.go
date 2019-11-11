package k8s

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleJobSpec represents an extremely simplified k8s job specification.
type SimpleJobSpec struct {
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	WorkingDir string            `json:"workingDir,omitempty"`
	Command    []string          `json:"command,omitempty"`
	CPU        resource.Quantity `json:"cpu,omitempty"`
	Memory     resource.Quantity `json:"memory,omitempty"`
	GPU        resource.Quantity `json:"gpu,omitempty"`
}

var (
	defaultTerminationMessagePolicy = corev1.TerminationMessageFallbackToLogsOnError

	// Do not restart failing jobs.
	defaultRestartPolicy = corev1.RestartPolicyOnFailure
	defaultBackoffLimit  = int32Ptr(0)
)

var defaultVolumes = []corev1.Volume{
	{
		Name: "storage",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "storage",
			},
		},
	},
}

var defaultVolumeMounts = []corev1.VolumeMount{
	{
		Name:      "storage",
		MountPath: "/storage",
	},
}

func (spec *SimpleJobSpec) resources() (*corev1.ResourceRequirements, error) {
	resources := &corev1.ResourceRequirements{Limits: corev1.ResourceList{}}

	if !spec.CPU.IsZero() {
		resources.Limits["cpu"] = spec.CPU
	}

	if !spec.Memory.IsZero() {
		resources.Limits["memory"] = spec.Memory
	}

	if !spec.GPU.IsZero() {
		resources.Limits["nvidia.com/gpu"] = spec.GPU
	}

	return resources, nil
}

func (spec *SimpleJobSpec) containers() ([]corev1.Container, error) {
	resources, err := spec.resources()
	if err != nil {
		return nil, err
	}

	// TODO: Implement "NVIDIA_XYZ" environment variables to fix e.g. `gpu: 0` problem.

	containers := []corev1.Container{
		{
			Name:         spec.Name,
			Image:        spec.Image,
			Command:      spec.Command,
			WorkingDir:   spec.WorkingDir,
			VolumeMounts: defaultVolumeMounts,
			Resources:    *resources,

			TerminationMessagePolicy: defaultTerminationMessagePolicy,
		},
	}

	return containers, nil
}

func (spec *SimpleJobSpec) meta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: spec.Name,
	}
}

// Expand expands the simplified job spec into a full job object.
func (spec *SimpleJobSpec) Expand() (*batchv1.Job, error) {
	containers, err := spec.containers()
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{
		ObjectMeta: spec.meta(),
		Spec: batchv1.JobSpec{
			BackoffLimit: defaultBackoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers:    containers,
					Volumes:       defaultVolumes,
					RestartPolicy: defaultRestartPolicy,
				},
			},
		},
	}

	return job, nil
}
