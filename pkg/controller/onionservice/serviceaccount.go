package onionservice

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

func (bc *OnionServiceController) reconcileServiceAccount(onionService *torv1alpha1.OnionService) error {
	serviceAccountName := onionService.ServiceAccountName()
	if serviceAccountName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("serviceAccount name must be specified"))
		return nil
	}

	serviceAccount, err := bc.KubernetesInformers.Core().V1().ServiceAccounts().Lister().ServiceAccounts(onionService.Namespace).Get(serviceAccountName)

	newServiceAccount := torServiceAccount(onionService)
	if errors.IsNotFound(err) {
		serviceAccount, err = bc.KubernetesClientSet.CoreV1().ServiceAccounts(onionService.Namespace).Create(newServiceAccount)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(serviceAccount, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, serviceAccount.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the service specs don't match, update
	if !serviceAccountEqual(serviceAccount, newServiceAccount) {
		serviceAccount, err = bc.KubernetesClientSet.CoreV1().ServiceAccounts(onionService.Namespace).Update(newServiceAccount)
	}

	if err != nil {
		return err
	}
	return nil
}

func serviceAccountEqual(a, b *corev1.ServiceAccount) bool {
	// TODO: actually detect differences
	return true
}

func torServiceAccount(onion *torv1alpha1.OnionService) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      onion.ServiceAccountName(),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   torv1alpha1.SchemeGroupVersion.Group,
					Version: torv1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
				}),
			},
		},
	}
}
