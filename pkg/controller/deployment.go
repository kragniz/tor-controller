package controller

import (
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
)

const (
	privateKeyVolume = "private-key"
	torConfigVolume  = "tor-config"
)

func (c *Controller) syncDeployment(onionService *v1alpha1.OnionService) error {
	deploymentName := fmt.Sprintf(deploymentNameFmt, onionService.Name)
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("deployment name must be specified"))
		return nil
	}

	// Get the deployment with the name specified in Foo.spec
	deployment, err := c.deploymentsLister.Deployments(onionService.Namespace).Get(deploymentName)

	// If the resource doesn't exist, we'll create it
	newDeployment := torDeployment(onionService)
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(onionService.Namespace).Create(newDeployment)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Foo resource, we should log
	// a warning to the event recorder and return.
	if !metav1.IsControlledBy(deployment, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the deployment specs don't match, update
	if !reflect.DeepEqual(deployment.Spec, newDeployment.Spec) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(onionService.Namespace).Update(newDeployment)
	}

	if err != nil {
		return err
	}
	return nil
}

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
							Command: []string{
								"tor",
							},
							Args: []string{
								"-f",
								"/etc/tor/tor-config",
							},
							ImagePullPolicy: "Never",

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      privateKeyVolume,
									MountPath: "/run/tor",
									ReadOnly:  true,
								},
								{
									Name:      torConfigVolume,
									MountPath: "/etc/",
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
									Items: []corev1.KeyToPath{
										{
											Key:  onion.Spec.PrivateKeySecret.Key,
											Path: "private_key",
										},
									},
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
									Items: []corev1.KeyToPath{
										{
											Key:  "tor-config",
											Path: "tor/tor-config",
											Mode: configPermissions(0700),
										},
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
