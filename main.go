package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type (
	MetricServer struct {
		Requests int
		Metrics  map[string]float64
	}
)

func (m *MetricServer) logRequest(r *http.Request) {
	m.Requests++
	log.Printf("%s - - \"%s %s\"", r.RemoteAddr, r.Method, r.RequestURI)
}

func (m *MetricServer) getRoot(w http.ResponseWriter, r *http.Request) {
	m.logRequest(r)
	_, err := io.WriteString(w, "<a href=\"/metrics\">Metrics</a>\n")
	if err != nil {
		panic(err)
	}
}

func (m *MetricServer) getMetrics(w http.ResponseWriter, r *http.Request) {
	m.logRequest(r)

	_, err := io.WriteString(w, fmt.Sprintf("http_requests_total: %f\n", float64(m.Requests)))
	if err != nil {
		panic(err)
	}

	for k, v := range m.Metrics {
		_, err := io.WriteString(w, fmt.Sprintf("%s: %f\n", k, v))
		if err != nil {
			panic(err)
		}
	}
}

func NewMetricServer(metricsFile string) *MetricServer {
	var server MetricServer

	server.Metrics = make(map[string]float64)

	log.Printf("reading metrics from %s", metricsFile)
	if fd, err := os.Open(metricsFile); err == nil {
		defer fd.Close()
		data, _ := ioutil.ReadAll(fd)
		err := json.Unmarshal(data, &server.Metrics)
		if err != nil {
			log.Printf("failed to read json from %s", metricsFile)
		}
	} else {
		log.Printf("failed to read metrics: %s", err)
	}

	log.Printf("found %d metrics", len(server.Metrics))

	return &server
}

func main() {
	var metricsFile string
	flag.StringVar(&metricsFile, "metrics", "metrics.json", "Metrics data")
	flag.Parse()

	server := NewMetricServer(metricsFile)

	http.HandleFunc("/", server.getRoot)
	http.HandleFunc("/metrics", server.getMetrics)

	err := http.ListenAndServe(":9283", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed")
	} else if err != nil {
		panic(err)
	}
}
