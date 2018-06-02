package v1alpha1

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "tor.k8s.io", Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&OnionService{},
		&OnionServiceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OnionServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OnionService `json:"items"`
}

// CRD Generation
func getFloat(f float64) *float64 {
	return &f
}

func getInt(i int64) *int64 {
	return &i
}

var (
	// Define CRDs for resources
	OnionServiceCRD = v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "onionservices.tor.k8s.io",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "tor.k8s.io",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "OnionService",
				Plural: "onionservices",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Type: "object",
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"extraConfig": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"ports": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"name": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"publicPort": v1beta1.JSONSchemaProps{
													Type:   "integer",
													Format: "int32",
												},
												"targetPort": v1beta1.JSONSchemaProps{
													Type:   "integer",
													Format: "int32",
												},
											},
										},
									},
								},
								"privateKeySecret": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"key": v1beta1.JSONSchemaProps{
											Type: "string",
										},
										"name": v1beta1.JSONSchemaProps{
											Type: "string",
										},
									},
								},
								"selector": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"version": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
									Enum:   []v1beta1.JSON{v1beta1.JSON{[]byte{48}}, v1beta1.JSON{[]byte{50}}, v1beta1.JSON{[]byte{51}}},
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"hostname": v1beta1.JSONSchemaProps{
									Type: "string",
								},
							},
						},
					},
				},
			},
		},
	}
)
