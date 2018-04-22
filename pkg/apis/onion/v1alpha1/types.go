/*
Copyright 2017 The Kubernetes Authors.

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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnionService is a specification for a OnionService resource
type OnionService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OnionServiceSpec   `json:"spec"`
	Status OnionServiceStatus `json:"status"`
}

// OnionServiceSpec is the spec for a OnionService resource
type OnionServiceSpec struct {
	Service OnionServiceBackend `json:"service"`

	PrivateKeySecret SecretReference `json:"privateKeySecret"`
}

type OnionServiceBackend struct {
	// Specifies the name of the referenced service.
	ServiceName string `json:"serviceName"`

	// Specifies the port of the referenced service.
	ServicePort intstr.IntOrString `json:"servicePort"`
}

// OnionServiceStatus is the status for a OnionService resource
type OnionServiceStatus struct {
	Hostname string `json:"hostname"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnionServiceList is a list of OnionService resources
type OnionServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []OnionService `json:"items"`
}

// SecretReference represents a Secret Reference
type SecretReference struct {
	// Name is unique within a namespace to reference a secret resource.
	Name string `json:"name,omitempty"`

	Key string `json:"key,omitempty"`
}
