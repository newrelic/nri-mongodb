FROM golang:1.18 as builder
COPY . /go/src/github.com/newrelic/nri-mongodb/
RUN cd /go/src/github.com/newrelic/nri-mongodb && \
    make && \
    strip ./bin/nri-mongodb

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-mongodb/bin/nri-mongodb /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-mongodb
USER 1000
