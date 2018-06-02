package inject

import (
	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	rscheme "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/scheme"
	"github.com/kragniz/tor-controller/pkg/controller/onionservice"
	"github.com/kragniz/tor-controller/pkg/inject/args"
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	rscheme.AddToScheme(scheme.Scheme)

	// Inject Informers
	Inject = append(Inject, func(arguments args.InjectArgs) error {
		Injector.ControllerManager = arguments.ControllerManager

		if err := arguments.ControllerManager.AddInformerProvider(&torv1alpha1.OnionService{}, arguments.Informers.Tor().V1alpha1().OnionServices()); err != nil {
			return err
		}

		// Add Kubernetes informers
		if err := arguments.ControllerManager.AddInformerProvider(&appsv1.Deployment{}, arguments.KubernetesInformers.Apps().V1().Deployments()); err != nil {
			return err
		}
		if err := arguments.ControllerManager.AddInformerProvider(&corev1.Service{}, arguments.KubernetesInformers.Core().V1().Services()); err != nil {
			return err
		}
		if err := arguments.ControllerManager.AddInformerProvider(&corev1.ConfigMap{}, arguments.KubernetesInformers.Core().V1().ConfigMaps()); err != nil {
			return err
		}
		if err := arguments.ControllerManager.AddInformerProvider(&corev1.Secret{}, arguments.KubernetesInformers.Core().V1().Secrets()); err != nil {
			return err
		}

		if c, err := onionservice.ProvideController(arguments); err != nil {
			return err
		} else {
			arguments.ControllerManager.AddController(c)
		}
		return nil
	})

	// Inject CRDs
	Injector.CRDs = append(Injector.CRDs, &torv1alpha1.OnionServiceCRD)
	// Inject PolicyRules
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{"tor.k8s.io"},
		Resources: []string{"*"},
		Verbs:     []string{"*"},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"apps",
		},
		Resources: []string{
			"deployments",
		},
		Verbs: []string{
			"create", "delete", "get", "list", "patch", "update", "watch",
		},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"",
		},
		Resources: []string{
			"services",
		},
		Verbs: []string{
			"create", "delete", "get", "list", "patch", "update", "watch",
		},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"",
		},
		Resources: []string{
			"configmaps",
		},
		Verbs: []string{
			"create", "delete", "get", "list", "patch", "update", "watch",
		},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"",
		},
		Resources: []string{
			"secrets",
		},
		Verbs: []string{
			"create", "delete", "get", "list", "patch", "update", "watch",
		},
	})
	// Inject GroupVersions
	Injector.GroupVersions = append(Injector.GroupVersions, schema.GroupVersion{
		Group:   "tor.k8s.io",
		Version: "v1alpha1",
	})
	Injector.RunFns = append(Injector.RunFns, func(arguments run.RunArguments) error {
		Injector.ControllerManager.RunInformersAndControllers(arguments)
		return nil
	})
}
