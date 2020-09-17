FROM registry.svc.ci.openshift.org/ocp/builder:rhel-8-golang-1.15-openshift-4.6 AS builder

ADD . /go/src/sigs.k8s.io/node-feature-discovery

WORKDIR /go/src/sigs.k8s.io/node-feature-discovery


RUN go install \
  -ldflags "-X sigs.k8s.io/node-feature-discovery/pkg/version.version=v0.4.0" \
  ./cmd/*
RUN install -D -m644 nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf

FROM registry.svc.ci.openshift.org/ocp/4.6:base

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /etc/kubernetes/node-feature-discovery /etc/kubernetes/node-feature-discovery
COPY --from=builder /go/bin/nfd-* /usr/bin/
