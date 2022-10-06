.PHONY: all test templates yamls
.FORCE:

GO_CMD ?= go
GO_FMT ?= gofmt

IMAGE_BUILD_CMD ?= podman build
IMAGE_BUILD_EXTRA_OPTS ?=
IMAGE_PUSH_CMD ?= podman push
CONTAINER_RUN_CMD ?= podman run

MDL ?= mdl

K8S_CODE_GENERATOR ?= ../code-generator

VERSION := $(shell git describe --tags --dirty --always)

IMAGE_REGISTRY ?= registry.k8s.io/nfd
IMAGE_TAG_NAME ?= $(VERSION)
IMAGE_EXTRA_TAG_NAMES ?=

IMAGE_NAME := node-feature-discovery
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG := $(IMAGE_REPO):$(IMAGE_TAG_NAME)
IMAGE_EXTRA_TAGS := $(foreach tag,$(IMAGE_EXTRA_TAG_NAMES),$(IMAGE_REPO):$(tag))

K8S_NAMESPACE ?= node-feature-discovery

OPENSHIFT ?=

# We use different mount prefix for local and container builds.
# Take CONTAINER_HOSTMOUNT_PREFIX from HOSTMOUNT_PREFIX if only the latter is specified
ifdef HOSTMOUNT_PREFIX
    CONTAINER_HOSTMOUNT_PREFIX := $(HOSTMOUNT_PREFIX)
else
    CONTAINER_HOSTMOUNT_PREFIX := /host-
endif
HOSTMOUNT_PREFIX ?= /host-

KUBECONFIG ?=
E2E_TEST_CONFIG ?=
E2E_PULL_IF_NOT_PRESENT ?= false

LDFLAGS = -ldflags "-s -w -X openshift/node-feature-discovery/pkg/version.version=$(VERSION) -X openshift/node-feature-discovery/pkg/utils/hostpath.pathPrefix=$(HOSTMOUNT_PREFIX)"

all: image

build:
	@mkdir -p bin
	$(GO_CMD) build -v -o bin $(LDFLAGS) ./cmd/...

install:
	$(GO_CMD) install -v $(LDFLAGS) ./cmd/...

local-image: image
image: yamls
	$(IMAGE_BUILD_CMD) --build-arg VERSION=$(VERSION) \
	    --build-arg HOSTMOUNT_PREFIX=$(CONTAINER_HOSTMOUNT_PREFIX) \
	    -t $(IMAGE_TAG) \
	    $(foreach tag,$(IMAGE_EXTRA_TAGS),-t $(tag)) \
	    $(IMAGE_BUILD_EXTRA_OPTS) ./

# clean NFD labels on all nodes
# devel only
deploy-prune:
	kubectl apply --validate=false -k deployment/overlays/prune/
	kubectl wait --for=condition=complete job -l app=nfd -n node-feature-discovery
	kubectl delete -k deployment/overlays/prune/

yamls:
	@./hack/kustomize.sh $(K8S_NAMESPACE) $(IMAGE_REPO) $(IMAGE_TAG_NAME)

deploy: yamls
	kubectl apply -k .

templates:
	@# Need to prepend each line in the sample config with spaces in order to
	@# fit correctly in the configmap spec.
	@sed s'/^/    /' deployment/components/worker-config/nfd-worker.conf.example > nfd-worker.conf.tmp
	@# The sed magic below replaces the block of text between the lines with start and end markers
	@start=NFD-WORKER-CONF-START-DO-NOT-REMOVE; \
	end=NFD-WORKER-CONF-END-DO-NOT-REMOVE; \
	sed -e "/$$start/,/$$end/{ /$$start/{ p; r nfd-worker.conf.tmp" \
	    -e "}; /$$end/p; d }" -i deployment/helm/node-feature-discovery/values.yaml
	@rm nfd-worker.conf.tmp

generate:
	go mod vendor
	go generate ./cmd/... ./pkg/... ./source/...
	rm -rf vendor/
	controller-gen object crd output:crd:stdout paths=./pkg/apis/... > deployment/base/nfd-crds/nodefeaturerule-crd.yaml
	cp deployment/base/nfd-crds/nodefeaturerule-crd.yaml deployment/helm/node-feature-discovery/manifests/
	rm -rf sigs.k8s.io
	$(K8S_CODE_GENERATOR)/generate-groups.sh client,informer,lister \
	    openshift/node-feature-discovery/pkg/generated \
	    openshift/node-feature-discovery/pkg/apis \
	    "nfd:v1alpha1" --output-base=. \
	    --go-header-file hack/boilerplate.go.txt
	rm -rf pkg/generated
	mv openshift/node-feature-discovery/pkg/generated pkg/
	rm -rf sigs.k8s.io

verify:	verify-gofmt

verify-gofmt:
ifeq (, $(GOFMT_CHECK))
	@echo "verify-gofmt: OK"
else
	@echo "verify-gofmt: ERROR: gofmt failed on the following files:"
	@echo "$(GOFMT_CHECK)"
	@echo ""
	@echo "For details, run: gofmt -d -s $(GOFMT_CHECK)"
	@echo ""
	@exit 1
endif

ci-lint:
	golangci-lint run --timeout 7m0s

lint:
	golint -set_exit_status ./...

mdlint:
	${CONTAINER_RUN_CMD} \
	--rm \
	--volume "${PWD}:/workdir:ro,z" \
	--workdir /workdir \
	ruby:slim \
	/workdir/scripts/test-infra/mdlint.sh

helm-lint:
	helm lint --strict deployment/helm/node-feature-discovery/

test:
	$(GO_CMD) test ./cmd/... ./pkg/... ./source/...

e2e-test:
	@if [ -z ${KUBECONFIG} ]; then echo "[ERR] KUBECONFIG missing, must be defined"; exit 1; fi
	$(GO_CMD) test -v ./test/e2e/ -args -nfd.repo=$(IMAGE_REPO) -nfd.tag=$(IMAGE_TAG_NAME) \
	    -kubeconfig=$(KUBECONFIG) \
	    -nfd.e2e-config=$(E2E_TEST_CONFIG) \
	    -nfd.pull-if-not-present=$(E2E_PULL_IF_NOT_PRESENT) \
	    -ginkgo.focus="\[kubernetes-sigs\]" \
	    $(if $(OPENSHIFT),-nfd.openshift,)
	$(GO_CMD) test -v ./test/e2e/ -args -nfd.repo=$(IMAGE_REPO) -nfd.tag=$(IMAGE_TAG_NAME)-minimal \
	    -kubeconfig=$(KUBECONFIG) \
	    -nfd.e2e-config=$(E2E_TEST_CONFIG) \
	    -nfd.pull-if-not-present=$(E2E_PULL_IF_NOT_PRESENT) \
	    -ginkgo.focus="\[kubernetes-sigs\]" \
	    $(if $(OPENSHIFT),-nfd.openshift,)

local-image-push: push
push:
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)
