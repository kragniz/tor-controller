package onionservice

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

func (bc *OnionServiceController) reconcileRolebinding(onionService *torv1alpha1.OnionService) error {
	roleName := onionService.RoleName()
	if roleName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("role name must be specified"))
		return nil
	}

	rolebinding, err := bc.KubernetesInformers.Rbac().V1().RoleBindings().Lister().RoleBindings(onionService.Namespace).Get(roleName)

	newRolebinding := torRolebinding(onionService)
	if errors.IsNotFound(err) {
		rolebinding, err = bc.KubernetesClientSet.RbacV1().RoleBindings(onionService.Namespace).Create(newRolebinding)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(rolebinding, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, rolebinding.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the service specs don't match, update
	if !rolebindingEqual(rolebinding, newRolebinding) {
		rolebinding, err = bc.KubernetesClientSet.RbacV1().RoleBindings(onionService.Namespace).Update(newRolebinding)
	}

	if err != nil {
		return err
	}
	return nil
}

func rolebindingEqual(a, b *rbacv1.RoleBinding) bool {
	// TODO: actually detect differences
	return true
}

func torRolebinding(onion *torv1alpha1.OnionService) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      onion.RoleName(),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   torv1alpha1.SchemeGroupVersion.Group,
					Version: torv1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
				}),
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.ServiceAccountKind,
				Name: onion.ServiceAccountName(),
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: onion.RoleName(),
		},
	}

}
