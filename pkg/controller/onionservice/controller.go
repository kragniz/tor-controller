package onionservice

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"
	appsv1 "k8s.io/api/apps/v1"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	torv1alpha1client "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/typed/tor/v1alpha1"
	torv1alpha1informer "github.com/kragniz/tor-controller/pkg/client/informers/externalversions/tor/v1alpha1"
	torv1alpha1lister "github.com/kragniz/tor-controller/pkg/client/listers/tor/v1alpha1"

	"github.com/kragniz/tor-controller/pkg/inject/args"
)

func (bc *OnionServiceController) Reconcile(k types.ReconcileKey) error {
	log.Printf("Implement the Reconcile function on onionservice.OnionServiceController to reconcile %s\n", k.Name)
	return nil
}

// +kubebuilder:controller:group=tor,version=v1alpha1,kind=OnionService,resource=onionservices
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:informers:group=apps,version=v1,kind=Deployment
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;watch;list
// +kubebuilder:informers:group=core,version=v1,kind=Service
type OnionServiceController struct {
	onionserviceLister torv1alpha1lister.OnionServiceLister
	onionserviceclient torv1alpha1client.TorV1alpha1Interface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	onionservicerecorder record.EventRecorder
}

func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &OnionServiceController{
		onionserviceLister: arguments.ControllerManager.GetInformerProvider(&torv1alpha1.OnionService{}).(torv1alpha1informer.OnionServiceInformer).Lister(),

		onionserviceclient:   arguments.Clientset.TorV1alpha1(),
		onionservicerecorder: arguments.CreateRecorder("OnionServiceController"),
	}

	// Create a new controller that will call OnionServiceController.Reconcile on changes to OnionServices
	gc := &controller.GenericController{
		Name:             "OnionServiceController",
		Reconcile:        bc.Reconcile,
		InformerRegistry: arguments.ControllerManager,
	}
	if err := gc.Watch(&torv1alpha1.OnionService{}); err != nil {
		return gc, err
	}

	if err := gc.WatchControllerOf(&appsv1.Deployment{}, eventhandlers.Path{bc.LookupFoo},
        predicates.ResourceVersionChanged); err != nil {
        return gc, err
}

	return gc, nil
}
