# multi stage build. needs docker >= 17.05
FROM golang:latest as builder
# Original MAINTAINER Nicolas Szalay <nico -at- rottenbytes -dot- info>

ADD . /go/src/github.com/ledgerhq/tzindex-exporter
WORKDIR /go/src/github.com/ledgerhq/tzindex-exporter
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -v -o tzindexExporter .

# final stage
FROM busybox
# Original MAINTAINER Nicolas Szalay <nico -at- rottenbytes -dot- info>

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ledgerhq/tzindex-exporter/tzindexExporter /bin/

ENTRYPOINT ["/bin/tzindexExporter"]
