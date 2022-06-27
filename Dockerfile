FROM golang:1.18 as builder
COPY . /go/src/github.com/newrelic/nri-memcached/
RUN cd /go/src/github.com/newrelic/nri-memcached && \
    make && \
    strip ./bin/nri-memcached

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-memcached/bin/nri-memcached /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-memcached
USER 1000
