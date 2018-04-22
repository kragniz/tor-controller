package controller

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
)

func torDeployment(onion *v1alpha1.OnionService) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "tor",
		"controller": onion.Name,
	}

	name := fmt.Sprintf(deploymentNameFmt, onion.Name)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Foo",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "tor",
							Image: "kragniz/kube-tor-daemon:latest",
						},
					},
				},
			},
		},
	}
}
