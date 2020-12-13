package resources

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/model/resources/v1alpha1"
	"k8ssphere.io/k8ssphere/pkg/model/resources/v1alpha1/pod"
)

type ResourceGetter struct {
	getters map[schema.GroupVersionResource]v1alpha1.Interface
}

func NewResourceGetter(factory informers.InformerFactory) *ResourceGetter {
	getters := make(map[schema.GroupVersionResource]v1alpha1.Interface)
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}] = pod.New(factory.KubernetesSharedInformerFactory())

	return &ResourceGetter{
		getters: getters,
	}
}

func (r *ResourceGetter) tryResource(resource string) v1alpha1.Interface {
	for k, v := range r.getters {
		if k.Resource == resource {
			return v
		}
	}
	return nil
}

func (r *ResourceGetter) Get(resource, namespace, name string) (runtime.Object, error) {
	getter := r.tryResource(resource)
	if getter == nil {
		return nil, nil
	}
	return getter.Get(namespace, name)
}
