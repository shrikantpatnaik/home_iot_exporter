package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Metrics struct {
	Metrics []Metric `json:"metrics"`
}

var (
	iotMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "iot",
		Name:      "metric",
		Help:      "Metric for Home IOT Network",
	},
		[]string{
			"name",
			"type",
		})
)

func init() {
	prometheus.MustRegister(iotMetric)
}

func metricPostHandler(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	_ = json.NewDecoder(r.Body).Decode(&metrics)
	for _, metric := range metrics.Metrics {
		iotMetric.WithLabelValues(metric.Name, metric.Type).Set(metric.Value)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	r.HandleFunc("/metrics", metricPostHandler).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		log.Println(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
