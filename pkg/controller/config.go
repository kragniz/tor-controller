package controller

import (
	"bytes"
	"text/template"

	"github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const configFormat = `
SocksPort 0
{{ range .HiddenServices }}
HiddenServiceDir {{ .ServiceDir }}
HiddenServicePort {{ .PublicPort }} {{ .ServiceClusterIP }}:{{ .ServicePort }}
{{ end }}
`

var configTemplate = template.Must(template.New("config").Parse(configFormat))

type onionService struct {
	ServiceName      string
	ServiceNamespace string
	ServiceClusterIP string
	ServiceDir       string
	ServicePort      int32
	PublicPort       int32
}

func buildTorConfig(onion *onionService) ([]byte, error) {
	var tmp bytes.Buffer
	err := configTemplate.Execute(&tmp, onion)
	if err != nil {
		return nil, err
	}
	return tmp.Bytes(), nil
}

func torConfigmap(onion *v1alpha1.OnionService, service *corev1.Service) (*corev1.ConfigMap, error) {
	s := onionService{
		ServiceName:      onion.Spec.Service.Name,
		ServiceNamespace: onion.Namespace,
		ServiceClusterIP: service.Spec.ClusterIP,
		ServiceDir:       "/run/tor/",
		ServicePort:      onion.Spec.Service.Port.IntVal,
		PublicPort:       onion.Spec.PublicPort.IntVal,
	}

	config, err := buildTorConfig(&s)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{}, nil
}
