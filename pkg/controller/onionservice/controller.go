package onionservice

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	"k8s.io/client-go/tools/record"

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
type OnionServiceController struct {
	onionserviceLister torv1alpha1lister.OnionServiceLister
	onionserviceclient torv1alpha1client.TorV1alpha1Interface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	onionservicerecorder record.EventRecorder
}

// ProvideController provides a controller that will be run at startup.  Kubebuilder will use codegeneration
// to automatically register this controller in the inject package
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	// INSERT INITIALIZATIONS FOR ADDITIONAL FIELDS HERE
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

	// IMPORTANT:
	// To watch additional resource types - such as those created by your controller - add gc.Watch* function calls here
	// Watch function calls will transform each object event into a OnionService Key to be reconciled by the controller.
	//
	// **********
	// For any new Watched types, you MUST add the appropriate // +kubebuilder:informer and // +kubebuilder:rbac
	// annotations to the OnionServiceController and run "kubebuilder generate.
	// This will generate the code to start the informers and create the RBAC rules needed for running in a cluster.
	// See:
	// https://godoc.org/github.com/kubernetes-sigs/kubebuilder/pkg/gen/controller#example-package
	// **********

	return gc, nil
}
