package app

import (
	"fmt"
	"github.com/spf13/cobra"
	k8sflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8ssphere.io/k8ssphere/cmd/ks-controller/app/options"
	"k8ssphere.io/k8ssphere/pkg/apis"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s"
	"k8ssphere.io/k8ssphere/pkg/utils/signals"
	"k8ssphere.io/k8ssphere/pkg/utils/term"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewControllerServer() *cobra.Command {
	option := options.NewControllerOptions()

	cmd := &cobra.Command{
		Use: "k8s-controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(option, signals.SetupSignalHandler())
		},
		SilenceUsage: true,
	}
	fs := cmd.Flags()
	nfs := option.Flags()
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

func Run(options *options.ControllerOptions, stopCh <-chan struct{}) error {
	client, err := k8s.NewKubernetes(options.KubernetesOptions)
	if err != nil {
		return nil
	}
	informerFactory := informers.NewInformerFactories(client.Kubernetes(),client.K8ssphere())

	mgrOptions := manager.Options{
		CertDir: options.WebhookCertDir,
		Port:    8443,
	}

	if options.LeaderElect {
		mgrOptions = manager.Options{
			CertDir:                 options.WebhookCertDir,
			Port:                    8443,
			LeaderElection:          options.LeaderElect,
			LeaderElectionNamespace: "k8ssphere-system",
			LeaderElectionID:        "k8s-controller-manager-leader-election",
			LeaseDuration:           &options.LeaderElection.LeaseDuration, //持有锁的时间
			RetryPeriod:             &options.LeaderElection.RetryPeriod,   //竞争获取锁的时间
			RenewDeadline:           &options.LeaderElection.RenewDeadline, //在更新租约的超时时间
		}
	}

	klog.V(0).Info("setting up manager")

	// Use 8443 instead of 443 cause we need root permission to bind port 443
	mgr, err := manager.New(client.Config(), mgrOptions)
	if err != nil {
		klog.Fatalf("unable to set up overall controller manager: %v", err)
	}
	if err = apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatalf("unable add APIs to scheme: %v", err)
	}
	//crd
	err = Execute(mgr, client, informerFactory, stopCh)
	if err != nil {
		klog.Fatalf("unable to register controllers to the manager: %v", err)
	}
	// Start cache data after all informer is registered
	klog.V(0).Info("Starting cache resource from apiserver...")
	informerFactory.Start(stopCh)

	klog.V(0).Info("Starting the controllers.")
	if err = mgr.Start(stopCh); err != nil {
		klog.Fatalf("unable to run the manager: %v", err)
	}

	return nil
}
