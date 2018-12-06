FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder
WORKDIR /go/src/sigs.k8s.io/node-feature-discovery
COPY . . 

ENV CMT_CAT_VERSION="v1.2.0"

RUN git clone --depth 1 -b $CMT_CAT_VERSION https://github.com/intel/intel-cmt-cat.git && \
    make -C intel-cmt-cat/lib install && \
    make -C rdt-discovery && \
    make -C rdt-discovery install

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN $GOPATH/bin/dep ensure

RUN go install \
  -ldflags "-s -w -X main.version=$NFD_VERSION" \
  sigs.k8s.io/node-feature-discovery
RUN install -D -m644 node-feature-discovery.conf.example /etc/kubernetes/node-feature-discovery/node-feature-discovery.conf

RUN go test .

# Create production image for running node feature discovery
FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base

COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/kubernetes/node-feature-discovery /etc/kubernetes/node-feature-discovery
RUN ldconfig
COPY --from=builder /go/bin/node-feature-discovery /usr/bin/node-feature-discovery

ENTRYPOINT ["/usr/bin/node-feature-discovery"]
