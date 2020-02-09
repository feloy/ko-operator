package controllers

import (
	"fmt"

	kov1alpha1 "github.com/feloy/ko-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createJob(kobuilder *kov1alpha1.KoBuilder, configName string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-job", kobuilder.Name),
			Namespace: kobuilder.Namespace,
		},
		Spec: batchv1.JobSpec{

			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: "ko-builder",
					RestartPolicy:      "Never",
					Containers: []corev1.Container{
						{
							Name:  "ko-builder",
							Image: "feloy/ko-builder:release-1.4.0",
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: configName,
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/etc/gcloud",
									Name:      "gcloud",
									ReadOnly:  true,
								},
								{
									MountPath: "/pod",
									Name:      "pod-info",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "gcloud",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "gcloud",
								},
							},
						},
						{
							Name: "pod-info",
							VolumeSource: corev1.VolumeSource{
								DownwardAPI: &corev1.DownwardAPIVolumeSource{
									Items: []corev1.DownwardAPIVolumeFile{
										{
											Path: "name",
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.name",
											},
										},
										{
											Path: "uid",
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.uid",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
