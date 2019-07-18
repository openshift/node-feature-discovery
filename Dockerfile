FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder

# Get (cache) deps in a separate layer
COPY go.mod go.sum /go/node-feature-discovery/

WORKDIR /go/node-feature-discovery

RUN go mod download

# Do actual build
COPY . /go/node-feature-discovery


RUN go install \
  -ldflags "-X sigs.k8s.io/node-feature-discovery/pkg/version.version=v0.4.0" \
  ./cmd/*
RUN install -D -m644 nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /etc/kubernetes/node-feature-discovery /etc/kubernetes/node-feature-discovery
COPY --from=builder /go/bin/nfd-* /usr/bin/
