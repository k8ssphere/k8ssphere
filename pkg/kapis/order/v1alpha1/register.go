package v1alpha1

import (
	restful "github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8ssphere.io/k8ssphere/pkg/apiserver/runtime"
	ksclient "k8ssphere.io/k8ssphere/pkg/client/clientset/versioned"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"net/http"
)

const (
	GroupName = "order.k8ssphere.io"
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}

/**

 */
func AddToContainer(c *restful.Container, factory informers.InformerFactory,kubernetes kubernetes.Interface,ksclient ksclient.Interface) error {
	webservice := runtime.NewWebService(GroupVersion)
	hander := NewHandler(factory,kubernetes,ksclient)
	webservice.Route(webservice.GET("create/{username}").
		To(hander.CreateOrder).
		Param(webservice.PathParameter("username", "the name of the project")).
		Returns(http.StatusOK, "OK", nil))

	c.Add(webservice)
	return nil
}
