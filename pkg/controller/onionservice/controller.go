package onionservice

import (
	"fmt"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/eventhandlers"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/predicates"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/record"

	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	torv1alpha1client "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/typed/tor/v1alpha1"
	torv1alpha1informer "github.com/kragniz/tor-controller/pkg/client/informers/externalversions/tor/v1alpha1"
	torv1alpha1lister "github.com/kragniz/tor-controller/pkg/client/listers/tor/v1alpha1"

	"github.com/kragniz/tor-controller/pkg/inject/args"
	"github.com/kragniz/tor-controller/pkg/onionaddr"
)

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

func (bc *OnionServiceController) Reconcile(k types.ReconcileKey) error {
	namespace, name := k.Namespace, k.Name
	onionService, err := bc.onionserviceLister.OnionServices(namespace).Get(name)
	if err != nil {
		// The OnionService resource may no longer exist, in which case we stop
		// processing.
		if apierrors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("onionService '%s' in work queue no longer exists", k))
			return nil
		}

		return err
	}

	err = bc.reconcileService(onionService)
	if err != nil {
		return err
	}

	err = bc.reconcileDeployment(onionService)
	if err != nil {
		return err
	}

	// Finally, we update the status block of the OnionService resource to reflect the
	// current state of the world
	err = bc.updateOnionServiceStatus(onionService)
	if err != nil {
		return err
	}

	bc.recorder.Event(onionService, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (bc *OnionServiceController) updateOnionServiceStatus(onionService *torv1alpha1.OnionService) error {
	onionServiceCopy := onionService.DeepCopy()

	privKeySecret, err := bc.KubernetesInformers.Core().V1().Secrets().Lister().Secrets(onionService.Namespace).Get(onionService.Spec.PrivateKeySecret.Name)
	if err != nil {
		return err
	}

	privKey := privKeySecret.Data[onionService.Spec.PrivateKeySecret.Key]
	hostname, err := onionaddr.GetAddress(privKey)
	if err != nil {
		return err
	}
	onionServiceCopy.Status.Hostname = hostname

	serviceName := onionService.ServiceName()
	service, err := bc.KubernetesInformers.Core().V1().Services().Lister().Services(onionService.Namespace).Get(serviceName)
	clusterIP := ""
	if errors.IsNotFound(err) {
		clusterIP = "0.0.0.0"
	} else if err != nil {
		return err
	} else {
		clusterIP = service.Spec.ClusterIP
	}

	onionServiceCopy.Status.TargetClusterIP = clusterIP

	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the Foo resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err = bc.Clientset.TorV1alpha1().OnionServices(onionService.Namespace).Update(onionServiceCopy)
	return err
}

// LookupOnionService looksup an OnionService from the lister
func (bc *OnionServiceController) LookupOnionService(r types.ReconcileKey) (interface{}, error) {
	return bc.Informers.Tor().V1alpha1().OnionServices().Lister().OnionServices(r.Namespace).Get(r.Name)
}

// +kubebuilder:controller:group=tor,version=v1alpha1,kind=OnionService,resource=onionservices
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:informers:group=apps,version=v1,kind=Deployment
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:informers:group=core,version=v1,kind=Service
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:informers:group=core,version=v1,kind=ConfigMap
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:informers:group=core,version=v1,kind=Secret
type OnionServiceController struct {
	args.InjectArgs

	onionserviceLister torv1alpha1lister.OnionServiceLister
	onionserviceclient torv1alpha1client.TorV1alpha1Interface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &OnionServiceController{
		InjectArgs: arguments,

		onionserviceLister: arguments.ControllerManager.GetInformerProvider(&torv1alpha1.OnionService{}).(torv1alpha1informer.OnionServiceInformer).Lister(),

		onionserviceclient: arguments.Clientset.TorV1alpha1(),
		recorder:           arguments.CreateRecorder("OnionServiceController"),
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

	if err := gc.WatchControllerOf(&appsv1.Deployment{}, eventhandlers.Path{bc.LookupOnionService},
		predicates.ResourceVersionChanged); err != nil {
		return gc, err
	}

	if err := gc.WatchControllerOf(&corev1.Service{}, eventhandlers.Path{bc.LookupOnionService},
		predicates.ResourceVersionChanged); err != nil {
		return gc, err
	}

	return gc, nil
}
