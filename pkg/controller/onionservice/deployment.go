package onionservice

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

const (
	privateKeyVolume = "private-key"
	torConfigVolume  = "tor-config"
)

func deploymentName(onion *torv1alpha1.OnionService) string {
	return fmt.Sprintf(deploymentNameFmt, onion.Name)
}

func (c *OnionServiceController) syncDeployment(onionService *torv1alpha1.OnionService) error {
	return nil
}

func deploymentEqual(a, b *appsv1.Deployment) bool {
	// TODO: actually detect differences
	return true
}

func torDeployment(onion *torv1alpha1.OnionService) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "tor",
		"controller": onion.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(onion),
			Namespace: onion.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(onion, schema.GroupVersionKind{
					Group:   torv1alpha1.SchemeGroupVersion.Group,
					Version: torv1alpha1.SchemeGroupVersion.Version,
					Kind:    "OnionService",
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
							Args: []string{
								"-f",
								"/etc/tor/tor-config",
							},
							ImagePullPolicy: "Never",

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      privateKeyVolume,
									MountPath: "/run/tor/service/private_key",
									SubPath:   onion.Spec.PrivateKeySecret.Key,
								},
								{
									Name:      torConfigVolume,
									MountPath: "/etc/tor/tor-config",
									SubPath:   "tor-config",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: privateKeyVolume,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: onion.Spec.PrivateKeySecret.Name,
								},
							},
						},
						{
							Name: torConfigVolume,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf(configmapNameFmt, onion.Name),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
