# tzindex-exporter
Prometheus exporter for the tezos indexer located at https://github.com/blockwatch-cc/tzindex

To build a linux binary : `go get ; CGO_ENABLED=0 GOOS=linux go build -v -o tzindexExporter .`
To run locally : `go run main.go`
To build a local binary : `go build`
