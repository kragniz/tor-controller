package onionservice

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

const (
	privateKeyVolume = "private-key"
	torConfigVolume  = "tor-config"
)

func (bc *OnionServiceController) reconcileDeployment(onionService *torv1alpha1.OnionService) error {
	deploymentName := onionService.DeploymentName()
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("%s/%s: deployment name must be specified", onionService.Namespace, onionService.Name))
		return nil
	}

	deployment, err := bc.KubernetesInformers.Apps().V1().Deployments().Lister().Deployments(onionService.Namespace).Get(deploymentName)

	// If the resource doesn't exist, we'll create it
	newDeployment := torDeployment(onionService)
	if apierrors.IsNotFound(err) {
		deployment, err = bc.KubernetesClientSet.AppsV1().Deployments(onionService.Namespace).Create(newDeployment)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Foo resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, onionService) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		bc.recorder.Event(onionService, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the deployment specs don't match, update
	if !deploymentEqual(deployment, newDeployment) {
		deployment, err = bc.KubernetesClientSet.AppsV1().Deployments(onionService.Namespace).Update(newDeployment)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	return nil
}

func deploymentEqual(a, b *appsv1.Deployment) bool {
	// TODO: actually detect differences
	return false
}

func torDeployment(onion *torv1alpha1.OnionService) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "tor",
		"controller": onion.Name,
	}

	privateKeyMountPath := "/run/tor/service/hs_ed25519_secret_key"
	if onion.Spec.GetVersion() == 2 {
		privateKeyMountPath = "/run/tor/service/private_key"
	}

	// allow not specifying a private key
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	if onion.Spec.PrivateKeySecret != (torv1alpha1.SecretReference{}) {
		volumes = []corev1.Volume{
			{
				Name: privateKeyVolume,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: onion.Spec.PrivateKeySecret.Name,
					},
				},
			},
		}

		volumeMounts = []corev1.VolumeMount{
			{
				Name:      privateKeyVolume,
				MountPath: privateKeyMountPath,
				SubPath:   onion.Spec.PrivateKeySecret.Key,
			},
		}
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      onion.DeploymentName(),
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
					ServiceAccountName: onion.ServiceAccountName(),
					Containers: []corev1.Container{
						{
							Name:  "tor",
							Image: "quay.io/kragniz/tor-daemon-manager:master",
							Args: []string{
								"-name",
								onion.Name,
								"-namespace",
								onion.Namespace,
							},
							ImagePullPolicy: "IfNotPresent",

							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}
