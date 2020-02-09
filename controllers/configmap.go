package controllers

import (
	"fmt"

	kov1alpha1 "github.com/feloy/ko-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createConfigMap(kobuilder *kov1alpha1.KoBuilder) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", kobuilder.Name),
			Namespace: kobuilder.Namespace,
		},
		Data: map[string]string{
			"REGISTRY":         kobuilder.Spec.Registry,
			"SERVICE_ACCOUNT":  kobuilder.Spec.ServiceAccount,
			"REPOSITORY":       kobuilder.Spec.Repository,
			"CHECKOUT":         kobuilder.Spec.Checkout,
			"CONFIG_PATH":      kobuilder.Spec.ConfigPath,
			"OWNER_APIVERSION": "kobuilders.ko.feloy.dev",
			"OWNER_CONTROLLER": "false",
			"OWNER_KIND":       "KoBuilder",
			"OWNER_NAME":       kobuilder.Name,
			"OWNER_UID":        string(kobuilder.UID),
		},
	}
}
