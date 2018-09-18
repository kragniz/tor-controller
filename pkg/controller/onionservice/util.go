package onionservice

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
)

func (r *ReconcileOnionService) updateOnionServiceStatus(onionService *torv1alpha1.OnionService) error {
	onionServiceCopy := onionService.DeepCopy()

	serviceName := onionService.ServiceName()
	service := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: serviceName, Namespace: onionService.Namespace}, service)
	clusterIP := ""
	if errors.IsNotFound(err) {
		clusterIP = "0.0.0.0"
	} else if err != nil {
		return err
	} else {
		clusterIP = service.Spec.ClusterIP
	}

	onionServiceCopy.Status.TargetClusterIP = clusterIP

	err = r.Update(context.TODO(), onionServiceCopy)
	return err
}
