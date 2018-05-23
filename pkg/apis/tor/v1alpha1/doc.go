// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/kragniz/tor-controller/pkg/apis/tor
// +k8s:defaulter-gen=TypeMeta
// +groupName=tor.k8s.io
package v1alpha1 // import "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
