#!/bin/sh -ex

go-to-protobuf \
   --output-base=. \
   --go-header-file ../../../../hack/boilerplate.go.txt \
   --proto-import ../../../../vendor/ \
   --packages +github.com/openshift/node-feature-discovery/pkg/apis/nfd/v1alpha1=v1alpha1 \
   --keep-gogoproto=false \
   --apimachinery-packages "-k8s.io/apimachinery/pkg/util/intstr"

mv github.com/openshift/node-feature-discovery/pkg/apis/nfd/v1alpha1/* .
rm -rf github.com/openshift
