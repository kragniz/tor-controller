package controller

import (
	"github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func torConfig(onion *v1alpha1.OnionService) (string, error) {
	return "", nil
}

func torConfigmap(onion *v1alpha1.OnionService) *corev1.ConfigMap {
	return &corev1.ConfigMap{}
}
