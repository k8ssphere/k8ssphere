package apiserver

import (
	"bytes"
	"context"
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"k8ssphere.io/k8ssphere/pkg/apiserver/config"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/kapis/order/v1alpha1"
	resources "k8ssphere.io/k8ssphere/pkg/kapis/resources/v1alpha1"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s"
	netutil "k8ssphere.io/k8ssphere/pkg/utils/net"
	"net/http"
	rt "runtime"
	"time"
)

const (
	// ApiRootPath defines the root path of all KubeSphere apis.
	ApiRootPath = "/kapis"
	// MimeMergePatchJson is the mime header used in merge request
	MimeMergePatchJson = "application/merge-patch+json"
	//
	MimeJsonPatchJson = "application/json-patch+json"
)

type APIServer struct {
	// number of kubesphere apiserver
	ServerCount int
	//
	Server *http.Server
	//
	Config *config.Config
	// webservice container, where all webservice defines
	container *restful.Container

	KubernetesClient k8s.Client

	InformerFactory informers.InformerFactory
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()
	s.container.Router(restful.CurlyRouter{})
	s.container.RecoverHandler(func(panicReason interface{}, httpWriter http.ResponseWriter) {
		logStackOnRecover(panicReason, httpWriter)
	})
	s.installRestfulApi()
	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}
	s.Server.Handler = s.container

	s.filterHandlerChain(stopCh)

	return nil
}

func (s APIServer) Run(stopCh <-chan struct{}) error {
	//
	err := s.waitForResourceSync(stopCh)
	cxt, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-stopCh
		s.Server.Shutdown(cxt)
	}()

	klog.V(0).Infof("Start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}
	return err
}

/**

 */
func (s *APIServer) filterHandlerChain(stopCh <-chan struct{}) {
	handler := s.Server.Handler

	//handler = filters.WithAuthorization(handler)

	//handler = filters.WithAuthentication(handler)

	s.Server.Handler = handler

}

func (s *APIServer) installRestfulApi() {
	urlruntime.Must(resources.AddToContainer(s.container, s.InformerFactory))
	urlruntime.Must(v1alpha1.AddToContainer(s.container, s.InformerFactory,s.KubernetesClient.Kubernetes(),s.KubernetesClient.K8ssphere()))

}

/**

 */
func (s *APIServer) waitForResourceSync(stopCh <-chan struct{}) error {
	discoveryClient := s.KubernetesClient.Kubernetes().Discovery()
	_, apiResourcesList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return err
	}

	isResourceExists := func(resource schema.GroupVersionResource) bool {
		for _, apiResource := range apiResourcesList {
			if apiResource.GroupVersion == resource.GroupVersion().String() {
				for _, rsc := range apiResource.APIResources {
					if rsc.Name == resource.Resource {
						return true
					}
				}
			}
		}
		return false
	}
	// resources we have to create informer first
	k8sGVRs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "nodes"},
		{Group: "", Version: "v1", Resource: "resourcequotas"},
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		{Group: "", Version: "v1", Resource: "secrets"},
		{Group: "", Version: "v1", Resource: "configmaps"},

		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},

		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "apps", Version: "v1", Resource: "replicasets"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "controllerrevisions"},

		{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},

		{Group: "batch", Version: "v1", Resource: "jobs"},
		{Group: "batch", Version: "v1beta1", Resource: "cronjobs"},

		{Group: "extensions", Version: "v1beta1", Resource: "ingresses"},

		{Group: "autoscaling", Version: "v2beta2", Resource: "horizontalpodautoscalers"},

		{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
	}

	for _, gvr := range k8sGVRs {
		if !isResourceExists(gvr) {
			klog.Warningf("resource %s not exists in the cluster", gvr)
		} else {
			_, err := s.InformerFactory.KubernetesSharedInformerFactory().ForResource(gvr)
			if err != nil {
				klog.Errorf("cannot create informer for %s", gvr)
				return err
			}
		}
	}

	s.InformerFactory.KubernetesSharedInformerFactory().Start(stopCh)
	s.InformerFactory.KubernetesSharedInformerFactory().WaitForCacheSync(stopCh)

	k8ssphereSharedInformer:= s.InformerFactory.K8ssphereSharedInformerFactory();

	ksGVRs := []schema.GroupVersionResource{
		{Group: "order.k8ssphere.io", Version: "v1alpha1", Resource: "cates"},
	}

	for _, gvr := range ksGVRs {
		if !isResourceExists(gvr) {
			klog.Warningf("resource %s not exists in the cluster", gvr)
		} else {
			_, err = k8ssphereSharedInformer.ForResource(gvr)
			if err != nil {
				return err
			}
		}
	}

	k8ssphereSharedInformer.Start(stopCh)
	k8ssphereSharedInformer.WaitForCacheSync(stopCh)

	return nil
}

/**

 */
func logRequestAndResponse(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	chain.ProcessFilter(req, resp)
	// Always log error response
	logWithVerbose := klog.V(4)
	if resp.StatusCode() > http.StatusBadRequest {
		logWithVerbose = klog.V(0)
	}

	logWithVerbose.Infof("%s - \"%s %s %s\" %d %d %dms",
		netutil.GetRequestIP(req.Request),
		req.Request.Method,
		req.Request.URL,
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
		time.Since(start)/time.Millisecond,
	)
}

/**

 */
func logStackOnRecover(panicReason interface{}, w http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("recover from panic situation: - %v\r\n", panicReason))
	for i := 2; ; i += 1 {
		_, file, line, ok := rt.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	klog.Errorln(buffer.String())

	headers := http.Header{}
	if ct := w.Header().Get("Content-Type"); len(ct) > 0 {
		headers.Set("Accept", ct)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
