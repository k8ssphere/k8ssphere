package order

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8ssphere.io/k8ssphere/pkg/apis/order/v1alpha1"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/model/resources"
	ksclient "k8ssphere.io/k8ssphere/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	CreateOrder(ctx context.Context, namespace string,cate * v1alpha1.Cate)(* v1alpha1.Cate,error)
}

type Order struct {
	kubernetes kubernetes.Interface
	ksclient ksclient.Interface
	resourceGetter *resources.ResourceGetter
}

func New(informers informers.InformerFactory, k8sclient kubernetes.Interface, ksclient ksclient.Interface) Interface {
	 return &Order{
	 	 kubernetes: k8sclient,
	 	 ksclient: ksclient,
	 	 resourceGetter: resources.NewResourceGetter(informers),
	 }
}

func (o Order)CreateOrder(ctx context.Context,namespace string,cate * v1alpha1.Cate) (* v1alpha1.Cate,error)  {
	  return o.ksclient.OrderV1alpha1().Cates(namespace).Create(ctx,cate,metav1.CreateOptions{})
}