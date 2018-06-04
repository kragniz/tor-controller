package local

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	client "github.com/kragniz/tor-controller/pkg/client/clientset/versioned"
	"github.com/kragniz/tor-controller/pkg/tordaemon"
)

var (
	namespace, onionServiceName string
)

func init() {
	flag.StringVar(&namespace, "namespace", "",
		"The namespace of the OnionService to manage.")

	flag.StringVar(&onionServiceName, "name", "",
		"The name of the OnionService to manage.")
}

type LocalManager struct {
	restConfig *rest.Config
	clientset  *client.Clientset

	stopCh chan struct{}

	daemon tordaemon.Tor

	// controller loop
	controller *Controller
}

func New(config *rest.Config) *LocalManager {
	t := &LocalManager{
		restConfig: config,
		stopCh:     make(chan struct{}),
		daemon:     tordaemon.Tor{},
	}
	return t
}

func (m *LocalManager) Run() error {
	var errors []error

	if onionServiceName == "" {
		errors = append(errors, fmt.Errorf("-name flag cannot be empty"))
	}
	if namespace == "" {
		errors = append(errors, fmt.Errorf("-namespace flag cannot be empty"))
	}
	if err := utilerrors.NewAggregate(errors); err != nil {
		return err
	}

	clientset, err := client.NewForConfig(m.restConfig)
	if err != nil {
		return err
	}
	m.clientset = clientset

	// listen to signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	m.signalHandler(signalCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.daemon.SetContext(ctx)

	os.Chmod("/run/tor/service", 0700)

	// start watching for API server events that trigger applies
	m.watchForNotifications()

	// Wait for all goroutines to exit
	<-m.stopCh

	return nil
}

func (m *LocalManager) Must(err error) *LocalManager {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return m
}

func (m *LocalManager) watchForNotifications() {

	// create the onionservice watcher
	onionListWatcher := cache.NewListWatchFromClient(
		m.clientset.TorV1alpha1().RESTClient(),
		"onionservices",
		namespace,
		fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", onionServiceName)),
	)

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(onionListWatcher, &v1alpha1.OnionService{}, time.Second*10, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.AddAfter(key, 2*time.Second)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.AddAfter(key, 2*time.Second)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.AddAfter(key, 2*time.Second)
			}
		},
	}, cache.Indexers{})

	m.controller = NewController(queue, indexer, informer, m)

	go m.controller.Run(1, m.stopCh)

}

func (m *LocalManager) signalHandler(ch chan os.Signal) {
	go func() {
		select {
		case <-m.stopCh:
			break
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				fmt.Println("received SIGHUP")

			case syscall.SIGINT:
				fmt.Println("received SIGINT")
				close(m.stopCh)

			case syscall.SIGTERM:
				fmt.Println("received SIGTERM")
				close(m.stopCh)
			}
		}
	}()
}
