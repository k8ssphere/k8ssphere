package options

import (
	"fmt"
	"github.com/spf13/cobra"
	k8sflag "k8s.io/component-base/cli/flag"
	"k8ssphere.io/k8ssphere/pkg/apiserver"
	"k8ssphere.io/k8ssphere/pkg/apiserver/config"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s"
	"k8ssphere.io/k8ssphere/pkg/utils/signals"
	"k8ssphere.io/k8ssphere/pkg/utils/term"
	"net/http"
)

type ServerRunConfig struct {
	ConfigFile string
	Config     *config.Config
	ListenIp   string
	Port       int
	SecurePort int
	TlsCert    string
	TlsKey     string
}

func NewServerRunConfig() *ServerRunConfig {
	s := ServerRunConfig{
		ConfigFile: "",
		Config:     config.NewConfig(),
		ListenIp:   "0.0.0.0",
		Port:       9090,
		SecurePort: 443,
		TlsCert:    "",
		TlsKey:     "",
	}
	return &s
}

/**

 */
func NewServerConfig() *cobra.Command {
	c := NewServerRunConfig()

	conf, err := config.TryLoadFromDisk()

	if err == nil {
		c = &ServerRunConfig{
			Config:     conf,
			ListenIp:   "0.0.0.0",
			Port:       9090,
			SecurePort: 443,
			TlsCert:    "",
			TlsKey:     "",
		}
	}
	cmd := &cobra.Command{
		Use: "ks-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(c, signals.SetupSignalHandler())
		},
		SilenceUsage: true,
	}
	fs := cmd.Flags()
	nfs := c.Flags()
	for _, f := range nfs.FlagSets {
		fs.AddFlagSet(f)
	}
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		k8sflag.PrintSections(cmd.OutOrStdout(), nfs, cols)
	})
	return cmd

}

/**
  new apiserver
*/
func Run(c *ServerRunConfig, stopCh <-chan struct{}) error {
	apiserver, err := c.NewApiServer(stopCh)
	if err != nil {
		return err
	}
	err = apiserver.PrepareRun(stopCh)
	if err != nil {
		return err
	}
	return apiserver.Run(stopCh)
}

/**

 */
func (s ServerRunConfig) NewApiServer(stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	apiserver := &apiserver.APIServer{
		Config: s.Config,
	}
	//init kubernetes
	kubernetesClient, err := k8s.NewKubernetes(s.Config.KubernetesOptions)
	if err != nil {
		return nil, err
	}
	apiserver.KubernetesClient = kubernetesClient


	apiserver.InformerFactory = informers.NewInformerFactories(kubernetesClient.Kubernetes(),kubernetesClient.K8ssphere())

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.Port),
	}

	apiserver.Server = server

	return apiserver, nil

}

/**

 */
func (s ServerRunConfig) Flags() (nfs k8sflag.NamedFlagSets) {
	s.Config.KubernetesOptions.AddFlags(nfs.FlagSet("kubernetes"), s.Config.KubernetesOptions)
	return nfs
}
