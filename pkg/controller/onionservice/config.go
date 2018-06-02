package onionservice

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

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

func configPermissions(p int32) *int32 {
	return &p
}

func buildTorConfig(onion *torv1alpha1.OnionService, serviceClusterIP string) (string, error) {
	ports := []portPair{}
	for _, p := range onion.Spec.Ports {
		port := portPair{
			ServicePort: p.TargetPort,
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
		Version:          onion.Spec.GetVersion(),
	}

	var tmp bytes.Buffer
	err := configTemplate.Execute(&tmp, s)
	if err != nil {
		return "", err
	}
	return tmp.String(), nil
}

func torConfigmap(onion *torv1alpha1.OnionService, serviceClusterIP string) (*corev1.ConfigMap, error) {
	config, err := buildTorConfig(onion, serviceClusterIP)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(onion),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   torv1alpha1.SchemeGroupVersion.Group,
					Version: torv1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
				}),
			},
		},
		Data: map[string]string{
			"tor-config": config,
		},
	}, nil
}

func configMapName(onion *torv1alpha1.OnionService) string {
	return fmt.Sprintf(configmapNameFmt, onion.Name)
}

func (bc *OnionServiceController) reconcileConfigmap(onionService *torv1alpha1.OnionService) error {
	configmapName := configMapName(onionService)
	if configmapName == "" {
		runtime.HandleError(fmt.Errorf("configmap name must be specified"))
		return nil
	}

	serviceName := serviceName(onionService)
	service, err := bc.KubernetesInformers.Core().V1().Services().Lister().Services(onionService.Namespace).Get(serviceName)
	clusterIP := ""
	if errors.IsNotFound(err) {
		clusterIP = "0.0.0.0"
	} else if err != nil {
		return err
	} else {
		clusterIP = service.Spec.ClusterIP
	}

	newConfigmap, err := torConfigmap(onionService, clusterIP)
	if err != nil {
		return err
	}

	configmap, err := bc.KubernetesInformers.Core().V1().ConfigMaps().Lister().ConfigMaps(onionService.Namespace).Get(configmapName)
	if errors.IsNotFound(err) {
		configmap, err = bc.KubernetesClientSet.CoreV1().ConfigMaps(onionService.Namespace).Create(newConfigmap)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(configmap, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, configmap.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	if !reflect.DeepEqual(configmap.Data, newConfigmap.Data) {
		configmap, err = bc.KubernetesClientSet.CoreV1().ConfigMaps(onionService.Namespace).Update(newConfigmap)
	}

	if err != nil {
		return err
	}

	return nil
}
