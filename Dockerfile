FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder

ADD . /go/src/sigs.k8s.io/node-feature-discovery

WORKDIR /go/src/sigs.k8s.io/node-feature-discovery


RUN go install \
  -ldflags "-s -w -X sigs.k8s.io/node-feature-discovery/pkg/version.version=v4.1" \
  ./cmd/*
RUN install -D -m644 node-feature-discovery.conf.example /etc/kubernetes/node-feature-discovery/node-feature-discovery.conf

RUN go test ./cmd/... ./test/unit/...

# Create production image for running node feature discovery
FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base

COPY --from=builder /etc/kubernetes/node-feature-discovery /etc/kubernetes/node-feature-discovery
COPY --from=builder /go/bin/nfd-* /usr/bin/

ENTRYPOINT ["/usr/bin/node-feature-discovery"]
