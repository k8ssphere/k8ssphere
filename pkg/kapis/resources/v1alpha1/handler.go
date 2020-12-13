package v1alpha1

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"k8ssphere.io/k8ssphere/pkg/model/resources"

	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/model/quotas"
)

type Handler struct {
	resourcesQuotes quotas.ResourcesQuotes
	resourceGetter  *resources.ResourceGetter
}

func NewHandler(factory informers.InformerFactory) *Handler {

	hander := &Handler{
		resourcesQuotes: quotas.NewResourcesQuotes(factory),
		resourceGetter:  resources.NewResourceGetter(factory),
	}
	return hander
}

func (h *Handler) handleGetResources(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	//resource := request.PathParameter("resources")
	//result, err := h.resourceGetter.Get("pods",namespace, "coredns-66bff467f8-c7582")
	result, err := h.resourcesQuotes.Get(namespace, "count/pods")
	fmt.Print(err)
	response.WriteEntity(result)
}
