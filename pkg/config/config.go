package config

import (
	"bytes"
	"text/template"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

const configFormat = `
SocksPort 0
HiddenServiceDir {{ .ServiceDir }}
HiddenServiceVersion {{ .Version }}
{{ range .Ports }}
HiddenServicePort {{ .PublicPort }} {{ $.ServiceClusterIP }}:{{ .ServicePort }}
{{ end }}
`

var configTemplate = template.Must(template.New("config").Parse(configFormat))

type onionService struct {
	ServiceName      string
	ServiceNamespace string
	ServiceClusterIP string
	ServiceDir       string
	Version          int
	Ports            []portPair
}

type portPair struct {
	ServicePort int32
	PublicPort  int32
}

func TorConfigForService(onion *torv1alpha1.OnionService) (string, error) {
	ports := []portPair{}
	for _, p := range onion.Spec.Ports {
		port := portPair{
			ServicePort: p.TargetPort,
			PublicPort:  p.PublicPort,
		}
		ports = append(ports, port)
	}

	s := onionService{
		ServiceName:      onion.ServiceName(),
		ServiceNamespace: onion.Namespace,
		ServiceClusterIP: onion.Status.TargetClusterIP,
		ServiceDir:       "/run/tor/service",
		Ports:            ports,
		Version:          onion.Spec.GetVersion(),
	}

	var tmp bytes.Buffer
	err := configTemplate.Execute(&tmp, s)
	if err != nil {
		return "", err
	}
	return tmp.String(), nil
}
