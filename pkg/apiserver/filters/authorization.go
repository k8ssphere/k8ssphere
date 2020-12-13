package filters

import (
	"k8s.io/klog/v2"
	"k8ssphere.io/k8ssphere/pkg/apiserver/request"
	"net/http"
)

/**

 */
func WithAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user, _ := request.UserFrom(ctx)
		klog.V(4).Infof("cans %s", user)

	})

}
