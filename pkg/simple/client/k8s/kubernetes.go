/**

 */
package k8s

import (
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v2/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s/options"
	k8ssphereclientset "k8ssphere.io/k8ssphere/pkg/client/clientset/versioned"
)

type Client interface {
	Kubernetes() kubernetes.Interface
	K8ssphere()  k8ssphereclientset.Interface
	Snapshot() snapshotclient.Interface
	ApiExtensions() apiextensionsclient.Interface
	Discovery() discovery.DiscoveryInterface
	Master() string
	Config() *rest.Config
}

type KubernetesClient struct {
	k8s           kubernetes.Interface
	k8ssphere	  k8ssphereclientset.Interface
	snapshot      snapshotclient.Interface
	apiextensions apiextensionsclient.Interface
	discovery     discovery.DiscoveryInterface
	master        string
	config        *rest.Config
}

func NewKubernetes(options *options.KubernetesOptions) (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, err
	}
	var kubernetesClient KubernetesClient
	kubernetesClient.k8s, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient.snapshot, err = snapshotclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient.apiextensions, err = apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient.discovery, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient.k8ssphere,err = k8ssphereclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient.config = config
	kubernetesClient.master = options.Master

	return &kubernetesClient, err
}

/**

 */
func (k KubernetesClient) Kubernetes() kubernetes.Interface {
	return k.k8s
}

func (k KubernetesClient) K8ssphere() k8ssphereclientset.Interface {
	return k.k8ssphere
}


func (k KubernetesClient) ApiExtensions() apiextensionsclient.Interface {
	return k.apiextensions
}

func (k KubernetesClient) Snapshot() snapshotclient.Interface {
	return k.snapshot
}

func (k KubernetesClient) Discovery() discovery.DiscoveryInterface {
	return k.discovery
}

func (k KubernetesClient) Master() string {
	return k.master
}

func (k KubernetesClient) Config() *rest.Config {
	return k.config
}
