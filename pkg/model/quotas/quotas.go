package quotas

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8ssphere.io/k8ssphere/pkg/informers"
)

const (
	podsKey                   = "count/pods"
	daemonsetsKey             = "count/daemonsets.apps"
	deploymentsKey            = "count/deployments.apps"
	ingressKey                = "count/ingresses.extensions"
	servicesKey               = "count/services"
	statefulsetsKey           = "count/statefulsets.apps"
	persistentvolumeclaimsKey = "persistentvolumeclaims"
	jobsKey                   = "count/jobs.batch"
	cronJobsKey               = "count/cronjobs.batch"
)

type ResourcesQuotes struct {
	informers informers.InformerFactory
}

func NewResourcesQuotes(informers informers.InformerFactory) ResourcesQuotes {
	return ResourcesQuotes{
		informers: informers,
	}
}

var supportedResources = map[string]schema.GroupVersionResource{
	deploymentsKey:            {Group: "apps", Version: "v1", Resource: "deployments"},
	daemonsetsKey:             {Group: "apps", Version: "v1", Resource: "daemonsets"},
	statefulsetsKey:           {Group: "apps", Version: "v1", Resource: "statefulsets"},
	podsKey:                   {Group: "", Version: "v1", Resource: "pods"},
	servicesKey:               {Group: "", Version: "v1", Resource: "services"},
	persistentvolumeclaimsKey: {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	ingressKey:                {Group: "extensions", Version: "v1beta1", Resource: "ingresses"},
	jobsKey:                   {Group: "batch", Version: "v1", Resource: "jobs"},
	cronJobsKey:               {Group: "batch", Version: "v1beta1", Resource: "cronjobs"},
}

func (res ResourcesQuotes) Get(namespace, resource string) (int, error) {

	genericInformer, _ := res.informers.KubernetesSharedInformerFactory().ForResource(supportedResources[resource])
	ret, _ := genericInformer.Lister().ByNamespace(namespace).List(labels.Everything())

	return len(ret), nil
}
