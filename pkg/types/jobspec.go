package types

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleJobSpec represents an extremely simplified k8s job specification.
type SimpleJobSpec struct {
	Name    string      `yaml:"name"`
	Image   string      `yaml:"image"`
	Command StringArray `yaml:"command"`
	GPU     uint        `yaml:"gpu"`
}

// Expand expands the simplified job spec into a full job object.
func (spec *SimpleJobSpec) Expand() *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.Name,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: Int32Ptr(0),
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

// Int32Ptr returns a pointer to the specified int32 value.
func Int32Ptr(i int32) *int32 { return &i }

// Int64Ptr returns a pointer to the specified int64 value.
func Int64Ptr(i int64) *int64 { return &i }
