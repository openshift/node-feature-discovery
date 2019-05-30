.PHONY: all test yamls
.FORCE:

GO_CMD := go
GO_FMT := gofmt

IMAGE_BUILD_CMD := podman build 

VERSION := v0.4.0

IMAGE_REGISTRY := quay.io/zvonkok
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

all: image

image:
	$(IMAGE_BUILD_CMD) -t $(IMAGE_TAG) ./
mock:
	mockery --name=FeatureSource --dir=source --inpkg --note="Re-generate by running 'make mock'"
	mockery --name=APIHelpers --dir=pkg/apihelper --inpkg --note="Re-generate by running 'make mock'"
	mockery --name=LabelerClient --dir=pkg/labeler --inpkg --note="Re-generate by running 'make mock'"

gofmt:
	@$(GO_FMT) -w -l $$(find . -name '*.go')

gofmt-verify:
	@out=`$(GO_FMT) -l -d $$(find . -name '*.go')`; \
	if [ -n "$$out" ]; then \
	    echo "$$out"; \
	    exit 1; \
	fi

ci-lint:
	golangci-lint run --timeout 5m0s

test:
	$(GO_CMD) test ./cmd/... ./pkg/...

e2e-test:
	$(GO_CMD) test -v ./test/e2e/ -args -nfd.repo=$(IMAGE_REPO) -nfd.tag=$(IMAGE_TAG_NAME) -kubeconfig=$(KUBECONFIG) -nfd.e2e-config=$(E2E_TEST_CONFIG)

push:
	$(IMAGE_PUSH_CMD) $(IMAGE_TAG)
