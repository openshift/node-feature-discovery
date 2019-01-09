.PHONY: all

IMAGE_BUILD_CMD := buildah bud 

QUAY_DOMAIN_NAME := quay.io
QUAY_REGISTRY_USER := zvonkok
DOCKER_IMAGE_NAME := node-feature-discovery

VERSION := v4.0

all: image

# To override QUAY_REGISTRY_USER use the -e option as follows:
# QUAY_REGISTRY_USER=<my-username> make docker -e.
image:
	$(IMAGE_BUILD_CMD) -t $(QUAY_DOMAIN_NAME)/$(QUAY_REGISTRY_USER)/$(DOCKER_IMAGE_NAME):$(VERSION) ./
