/***
Copyright 2019 Cisco Systems Inc. All rights reserved.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	aciawv1 "github.com/noironetworks/aci-containers/pkg/gbpcrd/apis/aci.aw/v1"
	versioned "github.com/noironetworks/aci-containers/pkg/gbpcrd/clientset/versioned"
	internalinterfaces "github.com/noironetworks/aci-containers/pkg/gbpcrd/informers/externalversions/internalinterfaces"
	v1 "github.com/noironetworks/aci-containers/pkg/gbpcrd/listers/aci.aw/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GBPServerInformer provides access to a shared informer and lister for
// GBPServers.
type GBPServerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GBPServerLister
}

type gBPServerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGBPServerInformer constructs a new informer for GBPServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGBPServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGBPServerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGBPServerInformer constructs a new informer for GBPServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGBPServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AciV1().GBPServers(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AciV1().GBPServers(namespace).Watch(options)
			},
		},
		&aciawv1.GBPServer{},
		resyncPeriod,
		indexers,
	)
}

func (f *gBPServerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGBPServerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gBPServerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&aciawv1.GBPServer{}, f.defaultInformer)
}

func (f *gBPServerInformer) Lister() v1.GBPServerLister {
	return v1.NewGBPServerLister(f.Informer().GetIndexer())
}
