/*


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

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	orderv1alpha1 "k8ssphere.io/k8ssphere/pkg/apis/order/v1alpha1"
	versioned "k8ssphere.io/k8ssphere/pkg/client/clientset/versioned"
	internalinterfaces "k8ssphere.io/k8ssphere/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "k8ssphere.io/k8ssphere/pkg/client/listers/order/v1alpha1"
)

// CateInformer provides access to a shared informer and lister for
// Cates.
type CateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.CateLister
}

type cateInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCateInformer constructs a new informer for Cate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCateInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCateInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCateInformer constructs a new informer for Cate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCateInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OrderV1alpha1().Cates(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OrderV1alpha1().Cates(namespace).Watch(context.TODO(), options)
			},
		},
		&orderv1alpha1.Cate{},
		resyncPeriod,
		indexers,
	)
}

func (f *cateInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCateInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cateInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&orderv1alpha1.Cate{}, f.defaultInformer)
}

func (f *cateInformer) Lister() v1alpha1.CateLister {
	return v1alpha1.NewCateLister(f.Informer().GetIndexer())
}
