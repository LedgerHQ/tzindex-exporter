# multi stage build. needs docker >= 17.05
FROM golang:latest as builder
MAINTAINER Nicolas Szalay <nico -at- rottenbytes -dot- info>

ADD . /go/src/github.com/rottenbytes/tzindex-exporter
WORKDIR /go/src/github.com/rottenbytes/tzindex-exporter
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -v -o tzindexExporter .

# final stage
FROM scratch
MAINTAINER Nicolas Szalay <nico -at- rottenbytes -dot- info>

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/rottenbytes/tzindex-exporter/tzindexExporter /

CMD ["/tzindexExporter"]
