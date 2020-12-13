package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	// Get retrieves a single object by its namespace and name
	Get(namespace, name string) (runtime.Object, error)
}
