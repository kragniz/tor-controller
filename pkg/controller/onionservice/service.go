package onionservice

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

func serviceName(onion *torv1alpha1.OnionService) string {
	return fmt.Sprintf(serviceNameFmt, onion.Name)
}

func (bc *OnionServiceController) syncService(onionService *torv1alpha1.OnionService) error {
	serviceName := serviceName(onionService)
	if serviceName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("service name must be specified"))
		return nil
	}

	service, err := bc.KubernetesInformers.Core().V1().Services().Lister().Services(onionService.Namespace).Get(serviceName)

	newService := torService(onionService)
	if errors.IsNotFound(err) {
		service, err = bc.KubernetesClientSet.CoreV1().Services(onionService.Namespace).Create(newService)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(service, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, service.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the service specs don't match, update
	if !serviceEqual(service, newService) {
		service, err = bc.KubernetesClientSet.CoreV1().Services(onionService.Namespace).Update(newService)
	}

	if err != nil {
		return err
	}
	return nil
}

func serviceEqual(a, b *corev1.Service) bool {
	// TODO: actually detect differences

	return true
}

func torService(onion *torv1alpha1.OnionService) *corev1.Service {
	ports := []corev1.ServicePort{}
	for _, p := range onion.Spec.Ports {
		port := corev1.ServicePort{
			Name:       p.Name,
			TargetPort: intstr.FromInt(int(p.TargetPort)),
			Port:       p.TargetPort,
		}
		ports = append(ports, port)
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(onion),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   torv1alpha1.SchemeGroupVersion.Group,
					Version: torv1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: onion.Spec.Selector,
			Ports:    ports,
		},
	}
}
