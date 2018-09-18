package local2

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kragniz/tor-controller/pkg/apis"
	"github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	"github.com/kragniz/tor-controller/pkg/config"
	"github.com/kragniz/tor-controller/pkg/tordaemon"
)

var (
	log = logf.Log.WithName("example-controller")

	namespace, onionServiceName string

	daemon tordaemon.Tor
)

func init() {
	flag.StringVar(&namespace, "namespace", "",
		"The namespace of the OnionService to manage.")

	flag.StringVar(&onionServiceName, "name", "",
		"The name of the OnionService to manage.")
}

func Run() {
	flag.Parse()
	logf.SetLogger(logf.ZapLogger(true))
	entryLog := log.WithName("entrypoint")

	var errs []error

	if onionServiceName == "" {
		errs = append(errs, fmt.Errorf("-name flag cannot be empty"))
	}
	if namespace == "" {
		errs = append(errs, fmt.Errorf("-namespace flag cannot be empty"))
	}
	if err := utilerrors.NewAggregate(errs); err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	daemon.SetContext(ctx)

	os.Chmod("/run/tor/service", 0700)

	// Setup a Manager
	mgr, err := manager.New(runtimeconfig.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		entryLog.Error(err, "unable to register schemes")
		os.Exit(1)
	}

	r := reconcileOnionService{client: mgr.GetClient(), log: log.WithName("reconciler")}
	c, err := controller.New("tor-local-controller", mgr, controller.Options{
		Reconciler: &r,
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch ReplicaSets and enqueue ReplicaSet object key
	if err := c.Watch(&source.Kind{Type: &v1alpha1.OnionService{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch OnionServices")
		os.Exit(1)
	}

	onionService := &v1alpha1.OnionService{}
	err = r.client.Get(
		context.TODO(),
		types.NamespacedName{Name: onionServiceName, Namespace: namespace},
		onionService,
	)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find existing OnionService")
	} else {
		err = r.syncOnionConfig(onionService)
		if err != nil {
			entryLog.Error(err, "unable to start tor")
			os.Exit(1)
		}
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

// reconcileReplicaSet reconciles ReplicaSets
type reconcileOnionService struct {
	client client.Client
	log    logr.Logger
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileOnionService{}

func (r *reconcileOnionService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convinient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", request)

	// Fetch the OnionService from the cache
	onionService := &v1alpha1.OnionService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, onionService)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find OnionService")
		return reconcile.Result{}, nil
	}

	// FIXME: find out how to filter the watch
	if onionService.Name != onionServiceName {
		return reconcile.Result{}, nil
	}

	if onionService.Namespace != namespace {
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch OnionService")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.syncOnionConfig(onionService)
}

func (r *reconcileOnionService) syncOnionConfig(onionService *v1alpha1.OnionService) error {
	torConfig, err := config.TorConfigForService(onionService)
	if err != nil {
		fmt.Printf("Generating config failed with %v\n", err)
		return err
	}

	reload := false

	torfile, err := ioutil.ReadFile("/run/tor/torfile")
	if os.IsNotExist(err) {
		reload = true
	} else if err != nil {
		return err
	}

	if string(torfile) != torConfig {
		reload = true
	}

	if reload {
		fmt.Printf("updating onion config for %s/%s\n", onionService.Namespace, onionService.Name)

		err = ioutil.WriteFile("/run/tor/torfile", []byte(torConfig), 0644)
		if err != nil {
			fmt.Printf("Writing config failed with %v\n", err)
			return err
		}

		daemon.Reload()
	}

	err = r.updateOnionServiceStatus(onionService)
	if err != nil {
		fmt.Printf("Updating status failed with %v\n", err)
		return err
	}

	return nil
}

func (r *reconcileOnionService) updateOnionServiceStatus(onionService *v1alpha1.OnionService) error {
	hostname, err := ioutil.ReadFile("/run/tor/service/hostname")
	if err != nil {
		fmt.Printf("Got this error when trying to find hostname: %v", err)
		hostname = []byte("")
	}

	newHostname := strings.TrimSpace(string(hostname))

	if newHostname != onionService.Status.Hostname {
		onionServiceCopy := onionService.DeepCopy()
		onionServiceCopy.Status.Hostname = newHostname

		err = r.client.Update(context.TODO(), onionServiceCopy)
		return err
	}
	return nil
}
