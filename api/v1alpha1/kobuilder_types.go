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

	// Foo is an example field of KoBuilder. Edit KoBuilder_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// KoBuilderStatus defines the observed state of KoBuilder
type KoBuilderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

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
