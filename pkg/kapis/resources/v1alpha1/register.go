package v1alpha1

import (
	restful "github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8ssphere.io/k8ssphere/pkg/apiserver/runtime"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"net/http"
)

const (
	GroupName = ""
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha3"}

func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}

/**

 */
func AddToContainer(c *restful.Container, factory informers.InformerFactory) error {
	webservice := runtime.NewWebService(GroupVersion)
	hander := NewHandler(factory)
	webservice.Route(webservice.GET("test/{namespace}/{resources}").
		To(hander.handleGetResources).
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Param(webservice.PathParameter("resources", "namespace level resource type, e.g. pods,jobs,configmaps,services.")).
		Returns(http.StatusOK, "OK", nil))

	c.Add(webservice)
	return nil
}
