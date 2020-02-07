/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kov1alpha1 "github.com/feloy/ko-operator/api/v1alpha1"
)

// KoBuilderReconciler reconciles a KoBuilder object
type KoBuilderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ko.feloy.dev,resources=kobuilders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ko.feloy.dev,resources=kobuilders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

func (r *KoBuilderReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()
	log := r.Log.WithValues("kobuilder", req.NamespacedName)

	kobuilder := new(kov1alpha1.KoBuilder)
	if err = r.Get(ctx, req.NamespacedName, kobuilder); err != nil {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		err = client.IgnoreNotFound(err)
		return
	}
	log.Info(fmt.Sprintf("kobuilder: %+v", kobuilder.Spec))

	var configName string
	if configName, err = r.applyConfig(ctx, log, kobuilder); err != nil {
		return
	}

	if err = r.applyKoBuilderJob(ctx, log, kobuilder, configName); err != nil {
		return
	}

	return
}

func (r *KoBuilderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kov1alpha1.KoBuilder{}).
		Complete(r)
}

func (r *KoBuilderReconciler) applyConfig(ctx context.Context, log logr.Logger, kobuilder *kov1alpha1.KoBuilder) (name string, err error) {

	config := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-config", kobuilder.Name),
			Namespace:    kobuilder.Namespace,
		},
		Data: map[string]string{
			"REGISTRY":        kobuilder.Spec.Registry,
			"SERVICE_ACCOUNT": kobuilder.Spec.ServiceAccount,
		},
	}

	controllerutil.SetControllerReference(kobuilder, config, r.Scheme)

	if err = r.Create(ctx, config); err != nil {
		log.Error(err, "unable to create configmap for kobuilder")
	}

	name = config.Name

	return
}

func (r *KoBuilderReconciler) applyKoBuilderJob(ctx context.Context, log logr.Logger, kobuilder *kov1alpha1.KoBuilder, configName string) (err error) {

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-job-", kobuilder.Name),
			Namespace:    kobuilder.Namespace,
		},
		Spec: batchv1.JobSpec{

			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: "ko-builder",
					RestartPolicy:      "Never",
					Containers: []corev1.Container{
						{
							Name:  "ko-builder",
							Image: "feloy/ko-builder:release-1.2.0",
							Env: []corev1.EnvVar{
								{
									Name:  "REPOSITORY",
									Value: kobuilder.Spec.Repository,
								},
								{
									Name:  "CHECKOUT",
									Value: kobuilder.Spec.Checkout,
								},
								{
									Name:  "CONFIG_PATH",
									Value: kobuilder.Spec.ConfigPath,
								},
							},
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

	controllerutil.SetControllerReference(kobuilder, job, r.Scheme)

	if err = r.Create(ctx, job); err != nil {
		log.Error(err, "unable to create job for kobuilder")
	}
	return
}
