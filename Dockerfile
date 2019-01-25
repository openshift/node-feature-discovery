FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder

ADD . /go/src/sigs.k8s.io/node-feature-discovery

WORKDIR /go/src/sigs.k8s.io/node-feature-discovery

ENV CMT_CAT_VERSION="v1.2.0"
ENV NFD_VERSION="v4.1"

RUN case $(uname -m) in \
        arm64) \
                echo "skip rdt on Arm64 platform" \
                ;; \
        *) \
                make -C intel-cmt-cat/lib install && \
                make -C rdt-discovery && \
                make -C rdt-discovery install \
                ;; \
        esac

RUN go install \
  -ldflags "-s -w -X sigs.k8s.io/node-feature-discovery/pkg/version.version=$NFD_VERSION" \
  ./cmd/*
RUN install -D -m644 node-feature-discovery.conf.example /etc/kubernetes/node-feature-discovery/node-feature-discovery.conf

RUN go test ./cmd/... ./test/unit/...

# Create production image for running node feature discovery
FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base

COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/kubernetes/node-feature-discovery /etc/kubernetes/node-feature-discovery
RUN ldconfig
COPY --from=builder /go/bin/nfd-* /usr/bin/

ENTRYPOINT ["/usr/bin/node-feature-discovery"]
