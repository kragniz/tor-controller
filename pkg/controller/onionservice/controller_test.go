package onionservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	. "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	. "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/typed/tor/v1alpha1"
)

var _ = Describe("OnionService controller", func() {
	var instance OnionService
	var privateKeySecret corev1.Secret
	var expectedKey types.ReconcileKey
	var client OnionServiceInterface
	var secretClient typedcorev1.SecretInterface

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
				PrivateKeySecret: SecretReference{
					Name: "test-private-key",
					Key:  "private_key",
				},
			},
		}
		instance.Name = "tor-1"
		expectedKey = types.ReconcileKey{
			Namespace: "default",
			Name:      "tor-1",
		}

		privateKeySecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-private-key",
				Namespace: "default",
			},
			StringData: map[string]string{
				"private_key": `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCYmaXyOrpMqh69/sSC0H5Cca7iK9h913VWjCcs4Ts9REeoftrA
tIOq/wRq4zwRxIYgxAM5PVPYLxpZXZC8hKi/sBp6deJudxTudV9Y9VPRnt4D4U7g
izQ5uFnQrEt5s/lG6qGOoWoLXoyvywbuhmnUG9VtTFUPmYXGUtJLIN4nAQIDAQAB
AoGABlQo+7jPzSopCDwurjYzZlIMrLigy+dJqIM7hLL6n/na9vP74o4Z/9d/uWcf
MKVz4pv+fjs65PZfI2GsBQWDcg9Yeg3HbG24clJg8HoCQ/Ctu0/N6nt1DNeDX6AC
LGUjTPK6krjbe6pppCVwB9JSsDAtcK+lVHymdFRPJ8TF6fkCQQDHnaDER6JaQP7u
2Nvcjdw+ejx03+73+FhtMzIxZh8qXQl+gOgrIyfJ9mhmw0T3S4lB7WAhj+BtSNDH
HDWLo/vHAkEAw7RJ8MJWFTXi06PLXveBcA5CGT6xZhF3wEyQOMuyYnUM0bTDhnUt
/3XTcK4pzu3VCTkowXfACLw1JOhfbqP29wJAdVwMgDnpjwy1lbG0Ggjhm238i255
Lhs5ygIWmYqD+kE26sRZO7twkkIoAXr+2jHz4enw4eqYNUhhTx8bsBzaUwJAGzgJ
DKZKyLps6NigIX41D3u8L7yrebG2QRWk/XE/RzhWZxhIFXxYwG4H0WU3xWMIvTao
93eLSuu6TH7RPxco8wJBAI5PmUNRc5qp4QUPGXYC3p6Q2yMDlKEIt6n/q+Pa6CS+
F+FK0Cv8YNH9QvZaCBYnbYyvRyU+jTz9XY3e67Vazs4=
-----END RSA PRIVATE KEY-----`,
			},
		}
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
		secretClient.Delete(privateKeySecret.Name, &metav1.DeleteOptions{})
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

			client = cs.TorV1alpha1().OnionServices("default")
			secretClient = ks.CoreV1().Secrets("default")

			_, err := secretClient.Create(&privateKeySecret)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for reconcile to happen
			Eventually(after, "10s", "100ms").Should(BeClosed())

			// INSERT YOUR CODE HERE - test conditions post reconcile
		})
	})
})
