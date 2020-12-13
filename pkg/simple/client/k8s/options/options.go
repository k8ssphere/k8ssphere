package options

import (
	"github.com/spf13/pflag"
	"k8ssphere.io/k8ssphere/pkg/utils/reflectutils"
	"os"
)

type KubernetesOptions struct {
	KubeConfig string `json:"kubeconfig" yaml:"kubeconfig"`
	// +optional
	Master string `json:"master,omitempty" yaml:"master"`
	// +optional
	QPS float32 `json:"qps,omitempty" yaml:"qps"`
	// +optional
	Burst int `json:"burst,omitempty" yaml:"burst"`
}

func NewKubernatesOptions() *KubernetesOptions {
	return &KubernetesOptions{
		KubeConfig: "",
		Master:     "",
		QPS:        1e6,
		Burst:      1e6,
	}
}

func (k *KubernetesOptions) ApplyTo(options *KubernetesOptions) {
	reflectutils.Override(options, k)
}

/**

 */
func (k *KubernetesOptions) Validate() []error {
	errors := []error{}
	if k.KubeConfig != "" {
		if _, err := os.Stat(k.KubeConfig); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

/**

 */
func (k *KubernetesOptions) AddFlags(fs *pflag.FlagSet, c *KubernetesOptions) {
	fs.StringVar(&k.KubeConfig, "kubeconfig", c.KubeConfig, ""+
		"Path for kubernetes kubeconfig file, if left blank, will use "+
		"in cluster way.")
	fs.StringVar(&k.Master, "master", c.Master, ""+
		"Used to generate kubeconfig for downloading, if not specified, will use host in kubeconfig.")
}
