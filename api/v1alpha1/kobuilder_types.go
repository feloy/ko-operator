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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KoBuilderSpec defines the desired state of KoBuilder
type KoBuilderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Registry is is the GCP registry used to pull built images
	Registry string `json:"registry,omitempty"`
	// ServiceAccount is the GCP service account having access to registry
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// Repository is the GitHub repository where the Go sources reside
	Repository string `json:"repository,omitempty"`
	// Checkout is the branch / commit / tag of the repository to checkout
	Checkout string `json:"checkout,omitempty"`
	// ConfigPath is the path in the repository containing the manifests to create Kubernetes resources
	ConfigPath string `json:"configPath,omitempty"`
}

// KoBuilderState is the state of the KoBuilder
type KoBuilderState string

const (
	// Deploying state when the job has been created and is not yet completed
	Deploying KoBuilderState = "Deploying"
	// Deployed state when the job has completed
	Deployed KoBuilderState = "Deployed"
	// ErrorDeploying state when the job has failed
	ErrorDeploying KoBuilderState = "ErrorDeploying"
	// Unknown state when the state is unknown
	Unknown KoBuilderState = "Unknown"
	// Updated state when the config has just been updated
	Updated KoBuilderState = "Updated"
)

// KoBuilderStatus defines the observed state of KoBuilder
type KoBuilderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// State indicates if the builder is "Deploying" or has "Deployed" the resources
	State KoBuilderState `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Repository",type=string,JSONPath=`.spec.repository`
// +kubebuilder:printcolumn:name="Checkout",type=string,JSONPath=`.spec.checkout`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`

// KoBuilder is the Schema for the kobuilders API
type KoBuilder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KoBuilderSpec   `json:"spec,omitempty"`
	Status KoBuilderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KoBuilderList contains a list of KoBuilder
type KoBuilderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KoBuilder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KoBuilder{}, &KoBuilderList{})
}
