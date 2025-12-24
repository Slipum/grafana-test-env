package main

import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "time"
)


var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "backend_http_requests_total",
            Help: "Number of HTTP requests",
        },
        []string{"path", "method"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "backend_http_request_duration_seconds",
            Help:    "Duration of HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"path"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

func main() {
    mux := http.NewServeMux()

    mux.Handle("/metrics", promhttp.Handler())

    mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues("/status"))
        defer timer.ObserveDuration()

        httpRequestsTotal.WithLabelValues("/status", r.Method).Inc()
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues("/api"))
        defer timer.ObserveDuration()

        httpRequestsTotal.WithLabelValues("/api", r.Method).Inc()
        // time.Sleep(100 * time.Millisecond) // имитация нагрузки
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "hello"}`))
    })

    http.ListenAndServe(":8080", mux)
}