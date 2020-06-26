.PHONY: all test yamls
.FORCE:

GO_CMD := go
GO_FMT := gofmt
GO111MODULE=on

IMAGE_BUILD_CMD := podman build
IMAGE_BUILD_EXTRA_OPTS :=
IMAGE_PUSH_CMD := podman push

VERSION := $(shell git describe --tags --dirty --always)

BIN := node-feature-discovery
IMAGE_REGISTRY := quay.io/openshift-psap
IMAGE_NAME := node-feature-discovery
IMAGE_TAG_NAME := $(VERSION)
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG := $(IMAGE_REPO):$(IMAGE_TAG_NAME)
K8S_NAMESPACE := kube-system
HOSTMOUNT_PREFIX := /host-
KUBECONFIG :=
E2E_TEST_CONFIG :=

yaml_templates := $(wildcard *.yaml.template)
yaml_instances := $(patsubst %.yaml.template,%.yaml,$(yaml_templates))

all: build

build:
	GOOS=linux $(GO_CMD) build -o $(BIN) -ldflags "-s -w -X sigs.k8s.io/node-feature-discovery/pkg/version.version=$NFD_VERSION -X sigs.k8s.io/node-feature-discovery/source.pathPrefix=$HOSTMOUNT_PREFIX" ./cmd/*

local-image: yamls
	$(IMAGE_BUILD_CMD) --build-arg NFD_VERSION=$(VERSION) \
		--build-arg HOSTMOUNT_PREFIX=$(HOSTMOUNT_PREFIX) \
		-t $(IMAGE_TAG) \
		$(IMAGE_BUILD_EXTRA_OPTS) ./

local-image-push:
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)

yamls: $(yaml_instances)

%.yaml: %.yaml.template .FORCE
	@echo "$@: namespace: ${K8S_NAMESPACE}"
	@echo "$@: image: ${IMAGE_TAG}"
	@sed -E \
	     -e s',^(\s*)name: node-feature-discovery # NFD namespace,\1name: ${K8S_NAMESPACE},' \
	     -e s',^(\s*)image:.+$$,\1image: ${IMAGE_TAG},' \
	     -e s',^(\s*)namespace:.+$$,\1namespace: ${K8S_NAMESPACE},' \
	     -e s',^(\s*)mountPath: "/host-,\1mountPath: "${HOSTMOUNT_PREFIX},' \
	     $< > $@

mock:
	mockery --name=FeatureSource --dir=source --inpkg --note="Re-generate by running 'make mock'"
	mockery --name=APIHelpers --dir=pkg/apihelper --inpkg --note="Re-generate by running 'make mock'"
	mockery --name=LabelerClient --dir=pkg/labeler --inpkg --note="Re-generate by running 'make mock'"

verify:	verify-gofmt

verify-gofmt:
	@out=`$(GO_FMT) -l -d $$(find . -name '*.go')`; \
	if [ -n "$$out" ]; then \
	    echo "$$out"; \
	    exit 1; \
	fi

gofmt:
	@$(GO_FMT) -w -l $$(find . -name '*.go')

ci-lint:
	golangci-lint run --timeout 5m0s

test:
	$(GO_CMD) test -x -v ./cmd/... ./pkg/...

test-e2e:
	$(GO_CMD) test -v ./test/e2e/ -args -nfd.repo=$(IMAGE_REPO) -nfd.tag=$(IMAGE_TAG_NAME) -kubeconfig=$(KUBECONFIG) -nfd.e2e-config=$(E2E_TEST_CONFIG)

clean:
		go clean
		rm -f $(BIN)

.PHONY: all build verify verify-gofmt clean local-image local-image-push test-e2e
