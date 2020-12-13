package v1alpha1

import (
	restful "github.com/emicklei/go-restful/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8ssphere.io/k8ssphere/pkg/apis/order/v1alpha1"
	ksclient "k8ssphere.io/k8ssphere/pkg/client/clientset/versioned"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/model/order"
)

type Handler struct {
	order order.Interface
}

func NewHandler(factory informers.InformerFactory,kubernetes kubernetes.Interface,ksclient ksclient.Interface) *Handler {

	hander := &Handler{
	    order: order.New(factory,kubernetes,ksclient),
	}
	return hander
}

func (h *Handler) CreateOrder(request *restful.Request, response *restful.Response) {
	username := request.PathParameter("username")
	cate := v1alpha1.Cate{
		   TypeMeta: metav1.TypeMeta {
                  	APIVersion:"order.k8ssphere.io/v1alpha1",
                  	Kind: "Cate",
				  },
		   ObjectMeta:metav1.ObjectMeta{
		   	        Name: "order-system",
				  },
	 	   Spec: v1alpha1.CateSpec{
	 	   	    Username: username,
		   },
	 }

	 order,_ := h.order.CreateOrder(request.Request.Context(),"default",&cate)
	response.WriteEntity(order)
}
