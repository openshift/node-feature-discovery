FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder

# Get (cache) deps in a separate layer
COPY go.mod go.sum /go/node-feature-discovery/

WORKDIR /go/node-feature-discovery

RUN go mod download

# Do actual build
COPY . /go/node-feature-discovery

ARG VERSION
ARG HOSTMOUNT_PREFIX

RUN go install \
  -ldflags "-s -w -X openshift/node-feature-discovery/pkg/version.version=$NFD_VERSION -X openshift/node-feature-discovery/source.pathPrefix=$HOSTMOUNT_PREFIX" \
  ./cmd/*
RUN install -D -m644 nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base

# Run as unprivileged user
USER 65534:65534

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /go/node-feature-discovery/nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf
COPY --from=builder /go/bin/* /usr/bin/
