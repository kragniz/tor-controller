package controller

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
)

const configFormat = `
SocksPort 0
HiddenServiceDir {{ .ServiceDir }}
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
	Ports            []portPair
}

type portPair struct {
	ServicePort int32
	PublicPort  int32
}

func configPermissions(p int32) *int32 {
	return &p
}

func buildTorConfig(onion *v1alpha1.OnionService, serviceClusterIP string) (string, error) {
	ports := []portPair{}
	for _, p := range onion.Spec.Ports {
		port := portPair{
			ServicePort: p.TargetPort.IntVal,
			PublicPort:  p.PublicPort,
		}
		ports = append(ports, port)
	}

	s := onionService{
		ServiceName:      serviceName(onion),
		ServiceNamespace: onion.Namespace,
		ServiceClusterIP: serviceClusterIP,
		ServiceDir:       "/run/tor/service",
		Ports:            ports,
	}

	var tmp bytes.Buffer
	err := configTemplate.Execute(&tmp, s)
	if err != nil {
		return "", err
	}
	return tmp.String(), nil
}

func torConfigmap(onion *v1alpha1.OnionService, serviceClusterIP string) (*corev1.ConfigMap, error) {
	config, err := buildTorConfig(onion, serviceClusterIP)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(configmapNameFmt, onion.Name),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
				}),
			},
		},
		Data: map[string]string{
			"tor-config": config,
		},
	}, nil
}

func (c *Controller) syncConfigmap(onionService *v1alpha1.OnionService) error {
	configmapName := fmt.Sprintf(configmapNameFmt, onionService.Name)
	if configmapName == "" {
		runtime.HandleError(fmt.Errorf("configmap name must be specified"))
		return nil
	}

	serviceName := serviceName(onionService)
	service, err := c.servicesLister.Services(onionService.Namespace).Get(serviceName)
	if err != nil {
		return err
	}

	clusterIP := service.Spec.ClusterIP

	newConfigmap, err := torConfigmap(onionService, clusterIP)
	if err != nil {
		return err
	}

	configmap, err := c.configmapsLister.ConfigMaps(onionService.Namespace).Get(configmapName)

	if errors.IsNotFound(err) {
		configmap, err = c.kubeclientset.CoreV1().ConfigMaps(onionService.Namespace).Create(newConfigmap)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(configmap, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, configmap.Name)
		c.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	if !reflect.DeepEqual(configmap.Data, newConfigmap.Data) {
		configmap, err = c.kubeclientset.CoreV1().ConfigMaps(onionService.Namespace).Update(newConfigmap)
	}

	if err != nil {
		return err
	}

	return nil
}
