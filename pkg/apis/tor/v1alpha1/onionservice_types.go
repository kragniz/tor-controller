package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OnionServiceSpec defines the desired state of OnionService
type OnionServiceSpec struct {
	// The list of ports that are exposed by this service.
	// +patchMergeKey=publicPort
	// +patchStrategy=merge
	Ports []ServicePort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"publicPort"`

	Selector map[string]string `json:"selector,omitempty"`

	PrivateKeySecret SecretReference `json:"privateKeySecret"`

	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Minimum=2
	Version int `json:"version"`

	ExtraConfig string `json:"extraConfig,omitempty"`
}

type ServicePort struct {
	// Optional if only one ServicePort is defined on this service.
	// +optional
	Name string `json:"name,omitempty"`

	// The port that will be exposed by this service.
	PublicPort int32 `json:"publicPort"`

	// Number or name of the port to access on the pods targeted by the service.
	// Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	// If this is a string, it will be looked up as a named port in the
	// target Pod's container ports. If this is not specified, the value
	// of the 'port' field is used (an identity map).
	// This field is ignored for services with clusterIP=None, and should be
	// omitted or set equal to the 'port' field.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
	// +optional
	TargetPort int32 `json:"targetPort,omitempty"`

	// TODO: figure out how to make kubebuilder allow IntOrString fields
	//TargetPort intstr.IntOrString `json:"targetPort,omitempty"`
}

// SecretReference represents a Secret Reference
type SecretReference struct {
	// Name is unique within a namespace to reference a secret resource.
	Name string `json:"name,omitempty"`

	Key string `json:"key,omitempty"`
}

// OnionServiceStatus defines the observed state of OnionService
type OnionServiceStatus struct {
	Hostname string `json:"hostname"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnionService
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=onionservices
type OnionService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OnionServiceSpec   `json:"spec,omitempty"`
	Status OnionServiceStatus `json:"status,omitempty"`
}
