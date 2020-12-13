/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filters

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
	"net/http"
)

// WithAuthentication installs authentication handler to handler chain.
// The following part is a little bit ugly, WithAuthentication also logs user failed login attempt
// if using basic auth. But only treats request with requestURI `/oauth/authorize` as login attempt
func WithAuthentication(handler http.Handler) http.Handler {

	s := serializer.NewCodecFactory(runtime.NewScheme()).WithoutConversion()
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		username, _, usingBasicAuth := req.BasicAuth()
		klog.Info("usename", username, usingBasicAuth, s)

		handler.ServeHTTP(w, req)
	})
}
