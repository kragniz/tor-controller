package inject

import (
	rbacv1 "k8s.io/api/rbac/v1"

	injectargs "github.com/kubernetes-sigs/kubebuilder/pkg/inject/args"
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"

	"github.com/kragniz/tor-controller/pkg/inject/args"
)

var (
	// Inject is used to add items to the Injector
	Inject []func(args.InjectArgs) error

	// Injector runs items
	Injector injectargs.Injector
)

// RunAll starts all of the informers and Controllers
func RunAll(rargs run.RunArguments, iargs args.InjectArgs) error {
	// Run functions to initialize injector
	for _, i := range Inject {
		if err := i(iargs); err != nil {
			return err
		}
	}

	if err := iargs.ControllerManager.AddInformerProvider(&rbacv1.Role{}, iargs.KubernetesInformers.Rbac().V1().Roles()); err != nil {
		return err
	}

	if err := iargs.ControllerManager.AddInformerProvider(&rbacv1.RoleBinding{}, iargs.KubernetesInformers.Rbac().V1().RoleBindings()); err != nil {
		return err
	}

	Injector.Run(rargs)
	<-rargs.Stop
	return nil
}
