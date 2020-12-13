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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "k8ssphere.io/k8ssphere/pkg/apis/order/v1alpha1"
)

// CateLister helps list Cates.
type CateLister interface {
	// List lists all Cates in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Cate, err error)
	// Cates returns an object that can list and get Cates.
	Cates(namespace string) CateNamespaceLister
	CateListerExpansion
}

// cateLister implements the CateLister interface.
type cateLister struct {
	indexer cache.Indexer
}

// NewCateLister returns a new CateLister.
func NewCateLister(indexer cache.Indexer) CateLister {
	return &cateLister{indexer: indexer}
}

// List lists all Cates in the indexer.
func (s *cateLister) List(selector labels.Selector) (ret []*v1alpha1.Cate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cate))
	})
	return ret, err
}

// Cates returns an object that can list and get Cates.
func (s *cateLister) Cates(namespace string) CateNamespaceLister {
	return cateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CateNamespaceLister helps list and get Cates.
type CateNamespaceLister interface {
	// List lists all Cates in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Cate, err error)
	// Get retrieves the Cate from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Cate, error)
	CateNamespaceListerExpansion
}

// cateNamespaceLister implements the CateNamespaceLister
// interface.
type cateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Cates in the indexer for a given namespace.
func (s cateNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Cate, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cate))
	})
	return ret, err
}

// Get retrieves the Cate from the indexer for a given namespace and name.
func (s cateNamespaceLister) Get(name string) (*v1alpha1.Cate, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("cate"), name)
	}
	return obj.(*v1alpha1.Cate), nil
}
