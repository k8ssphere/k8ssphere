package app

import (
	"k8s.io/klog/v2"
	"k8ssphere.io/k8ssphere/pkg/controller/order"
	"k8ssphere.io/k8ssphere/pkg/informers"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Execute(mgr manager.Manager, client k8s.Client, factory informers.InformerFactory, stopCh <-chan struct{}) error {

	k8ssphereInformer := factory.K8ssphereSharedInformerFactory();
	var orderController manager.Runnable

	orderController = order.NewOrderController(client.Kubernetes(),client.K8ssphere(),k8ssphereInformer.Order().V1alpha1().Cates())

	controllers := map[string]manager.Runnable{
		order.ControllerAgentName :orderController,
	}

	for name, controller := range controllers {
		if controller == nil {
			klog.V(4).Infof("%s is not going to run due to dependent component disabled.", name)
			continue
		}

		if err := mgr.Add(controller); err != nil {
			klog.Error(err, "add controller to manager failed", "name", name)
			return err
		}
	}

	return nil
}
