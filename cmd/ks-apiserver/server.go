package main

import (
	"k8s.io/klog/v2"
	"k8ssphere.io/k8ssphere/cmd/ks-apiserver/options"
)

/**

 */
func main() {
	//run init
	cmd := options.NewServerConfig()
	error := cmd.Execute()
	if error != nil {
		klog.Info(error)
	}
}
