FROM registry.svc.ci.openshift.org/ocp/builder:rhel-8-golang-1.15-openshift-4.7 AS builder

WORKDIR /go/node-feature-discovery

# Do actual build
COPY . /go/node-feature-discovery

RUN go install \
  -ldflags "-X sigs.k8s.io/node-feature-discovery/pkg/version.version=v0.6.0" \
  ./cmd/*
RUN install -D -m644 nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf

FROM registry.svc.ci.openshift.org/ocp/4.7:base

# Run as unprivileged user
USER 65534:65534

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /go/node-feature-discovery/nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf
COPY --from=builder /go/bin/* /usr/bin/
