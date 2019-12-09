package main

import (
	"encoding/json"
	//"flag"
	"github.com/namsral/flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
Response is :
{
  "status": "synced",
  "blocks": 715655,
  "indexed": 715655,
  "progress": 1
}

*/

type indexerStatus struct {
	Status string `json:"status"`
  Blocks int64 `json:"blocks"`
  Indexed int64 `json:"indexed"`
  Progress float64 `json:"progress"`
}

func getIndexerStatus(url string) {
	go func() {
		for {
			resp, err := http.Get(url)

			if err != nil {
				log.Println(err)
				time.Sleep(30 * time.Second)
				continue
			}

			body, _ := ioutil.ReadAll(resp.Body)
			var status indexerStatus
			jsonErr := json.Unmarshal(body, &status)
			if jsonErr != nil {
				log.Println(jsonErr)
			}

			blockHeightSeen.Set(float64(status.Blocks))
			blockHeightIndexed.Set(float64(status.Indexed))
			progress.Set(float64(status.Progress))
			// Indexer status (connecting, syncing, synced, failed). we care only about "synced"
			if status.Status == "synced" {
				synced.Set(float64(1))
			} else {
				synced.Set(float64(0))
			}

			// no more that one call every 30 seconds
			time.Sleep(30 * time.Second)
		}
	}()
}

// From https://tzstats.com/docs/api/index.html#indexer-status
var (
	blockHeightSeen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "block_height_seen",
		Help: "Seen height of the blockchain",
	})

	blockHeightIndexed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "block_height_indexed",
		Help: "Height of the current indexed block on this node",
	})

	synced = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "synced",
		Help: "Synced status returned by the indexer",
	})

	progress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "progress",
		Help: "progress status returned by the indexer. Percent, from 0 to 1",
	})
)

func main() {
	urlPtr := flag.String("url", "https://127.0.0.1/explorer/status", "URL for tzindexer status")
	portPtr := flag.String("port", "7392", "Port to expose")

	flag.Parse()

	fmt.Println("Will call tzindexer at", *urlPtr)
	fmt.Println("Listening on port", *portPtr)

	getIndexerStatus(*urlPtr)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*portPtr, nil)
}
