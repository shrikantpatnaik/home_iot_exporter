package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			iotMetric.WithLabelValues("hello", "world").Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	iotMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sp_home",
		Subsystem: "iot",
		Name:      "homeiot_metric",
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

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
