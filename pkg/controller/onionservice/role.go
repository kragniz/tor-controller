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

func (bc *OnionServiceController) reconcileRole(onionService *torv1alpha1.OnionService) error {
	roleName := onionService.RoleName()
	if roleName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("role name must be specified"))
		return nil
	}

	role, err := bc.KubernetesInformers.Rbac().V1().Roles().Lister().Roles(onionService.Namespace).Get(roleName)

	newRole := torRole(onionService)
	if errors.IsNotFound(err) {
		role, err = bc.KubernetesClientSet.RbacV1().Roles(onionService.Namespace).Create(newRole)
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(role, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, role.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the service specs don't match, update
	if !roleEqual(role, newRole) {
		role, err = bc.KubernetesClientSet.RbacV1().Roles(onionService.Namespace).Update(newRole)
	}

	if err != nil {
		return err
	}
	return nil
}

func roleEqual(a, b *rbacv1.Role) bool {
	// TODO: actually detect differences
	return true
}

func torRole(onion *torv1alpha1.OnionService) *rbacv1.Role {
	return &rbacv1.Role{
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
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{torv1alpha1.SchemeGroupVersion.Group},
				Verbs:     []string{"get", "list", "watch", "update", "patch"},
				Resources: []string{
					"onionservices",
				},
			},
			{
				APIGroups: []string{""},
				Verbs:     []string{"create", "update", "patch"},
				Resources: []string{"events"},
			},
		},
	}
}
