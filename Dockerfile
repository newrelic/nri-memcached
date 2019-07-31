FROM golang:1.10 as builder
RUN go get -d github.com/newrelic/nri-memcached/... && \
    cd /go/src/github.com/newrelic/nri-memcached && \
    make && \
    strip ./bin/nr-memcached

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-memcached/bin/nr-memcached /var/db/newrelic-infra/newrelic-integrations/bin/nr-memcached
COPY --from=builder /go/src/github.com/newrelic/nri-memcached/memcached-definition.yml /var/db/newrelic-infra/newrelic-integrations/definition.yml