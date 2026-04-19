package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const namespace = "artisancode_backend"

type Metrics struct {
	registry          *prometheus.Registry
	requestsTotal     *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
	requestSizeBytes  *prometheus.HistogramVec
	responseSizeBytes *prometheus.HistogramVec
	requestsInFlight  *prometheus.GaugeVec
	buildInfo         *prometheus.GaugeVec
}

func New(appName string, appVersion string, appEnvironment string) *Metrics {
	registry := prometheus.NewRegistry()

	requestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests processed by the backend.",
	}, []string{"method", "route", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency distributions in seconds.",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method", "route", "status"})

	requestSizeBytes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "request_size_bytes",
		Help:      "Approximate size of HTTP requests in bytes.",
		Buckets:   prometheus.ExponentialBuckets(256, 2, 10),
	}, []string{"method", "route"})

	responseSizeBytes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "response_size_bytes",
		Help:      "Size of HTTP responses in bytes.",
		Buckets:   prometheus.ExponentialBuckets(256, 2, 10),
	}, []string{"method", "route", "status"})

	requestsInFlight := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Current number of in-flight HTTP requests.",
	}, []string{"method", "route"})

	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "build_info",
		Help:      "Build information for the running backend service.",
	}, []string{"app_name", "version", "environment"})

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requestsTotal,
		requestDuration,
		requestSizeBytes,
		responseSizeBytes,
		requestsInFlight,
		buildInfo,
	)

	buildInfo.WithLabelValues(sanitizeLabelValue(appName), sanitizeLabelValue(appVersion), sanitizeLabelValue(appEnvironment)).Set(1)

	return &Metrics{
		registry:          registry,
		requestsTotal:     requestsTotal,
		requestDuration:   requestDuration,
		requestSizeBytes:  requestSizeBytes,
		responseSizeBytes: responseSizeBytes,
		requestsInFlight:  requestsInFlight,
		buildInfo:         buildInfo,
	}
}

func sanitizeLabelValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "unknown"
	}

	return trimmed
}
