/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	time "time"

	onion_v1alpha1 "github.com/kragniz/kube-onions/pkg/apis/onion/v1alpha1"
	versioned "github.com/kragniz/kube-onions/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kragniz/kube-onions/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kragniz/kube-onions/pkg/client/listers/onion/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// OnionServiceInformer provides access to a shared informer and lister for
// OnionServices.
type OnionServiceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.OnionServiceLister
}

type onionServiceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewOnionServiceInformer constructs a new informer for OnionService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewOnionServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredOnionServiceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredOnionServiceInformer constructs a new informer for OnionService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredOnionServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OnionV1alpha1().OnionServices(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OnionV1alpha1().OnionServices(namespace).Watch(options)
			},
		},
		&onion_v1alpha1.OnionService{},
		resyncPeriod,
		indexers,
	)
}

func (f *onionServiceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredOnionServiceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *onionServiceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&onion_v1alpha1.OnionService{}, f.defaultInformer)
}

func (f *onionServiceInformer) Lister() v1alpha1.OnionServiceLister {
	return v1alpha1.NewOnionServiceLister(f.Informer().GetIndexer())
}
