FROM golang:1.9 as builder
RUN go get -d github.com/newrelic/nri-mongodb/... && \
    cd /go/src/github.com/newrelic/nri-mongodb && \
    make && \
    strip ./bin/nr-mongodb

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-mongodb/bin/nr-mongodb /var/db/newrelic-infra/newrelic-integrations/bin/nr-mongodb
COPY --from=builder /go/src/github.com/newrelic/nri-mongodb/mongodb-definition.yml /var/db/newrelic-infra/newrelic-integrations/definition.yml
