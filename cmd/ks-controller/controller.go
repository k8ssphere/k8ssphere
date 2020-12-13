package main

import (
	"k8ssphere.io/k8ssphere/cmd/ks-controller/app"
	"os"
)

func main() {
	cmd := app.NewControllerServer()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
