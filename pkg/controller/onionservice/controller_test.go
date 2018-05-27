package onionservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	. "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/typed/tor/v1alpha1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic tests

var _ = Describe("OnionService controller", func() {
	var instance OnionService
	var expectedKey types.ReconcileKey
	var client OnionServiceInterface

	BeforeEach(func() {
		instance = OnionService{
			Spec: OnionServiceSpec{
				Ports: []ServicePort{
					ServicePort{
						Name:       "port1",
						PublicPort: 80,
						TargetPort: 8080,
					},
				},
				Selector: map[string]string{
					"app": "test",
				},
			},
		}
		instance.Name = "tor-1"
		expectedKey = types.ReconcileKey{
			Namespace: "default",
			Name:      "tor-1",
		}
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			after := make(chan struct{})
			ctrl.AfterReconcile = func(key types.ReconcileKey, err error) {
				defer func() {
					// Recover in case the key is reconciled multiple times
					defer func() { recover() }()
					close(after)
				}()
				defer GinkgoRecover()
				Expect(key).To(Equal(expectedKey))
				Expect(err).ToNot(HaveOccurred())
			}

			// Create the instance

			client = cs.TorV1alpha1().OnionServices("default")

			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for reconcile to happen
			Eventually(after, "10s", "100ms").Should(BeClosed())

			// INSERT YOUR CODE HERE - test conditions post reconcile
		})
	})
})
