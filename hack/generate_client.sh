#!/bin/bash

set -e

GV="order:v1alpha1"

rm -rf ./pkg/client
./hack/generate_group.sh "client,lister,informer" k8ssphere.io/k8ssphere/pkg/client k8ssphere.io/k8ssphere/pkg/apis "$GV" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv k8ssphere.io/k8ssphere/pkg/client ./pkg/
rm -rf ./k8ssphere.io
