package pod

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8ssphere.io/k8ssphere/pkg/model/resources/v1alpha1"
)

type podsGetter struct {
	informer informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha1.Interface {
	return &podsGetter{informer: sharedInformers}
}

func (p podsGetter) Get(namespace, name string) (runtime.Object, error) {
	return p.informer.Core().V1().Pods().Lister().Pods(namespace).Get(name)
}
