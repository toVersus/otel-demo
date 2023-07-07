package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	meterName     = "github.com/toVersus/otel-demo"
	metricsPrefix = "oteldemo"
)

var (
	httpInFlightRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"handler", "code", "method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of all HTTP requests",
	}, []string{"handler", "code", "method"})

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)
)

func InitMetricProvider() {
	exporter, err := promexporter.New()
	if err != nil {
		log.Fatalf("failed to initialize prometheus exporter: %s", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(newResource()),
		metric.WithReader(exporter),
	)

	provider.Meter(meterName)
	otel.SetMeterProvider(provider)
}

func SetupMetricsExporter() *http.Server {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		httpInFlightRequests,
		httpRequestsTotal,
		httpRequestDuration,
		responseSize,
	)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	return &http.Server{
		Addr:    ":2223",
		Handler: router,
	}
}

func ServeMetrics(srv *http.Server) {
	log.Printf("Metrics service running at :2223")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("failed to run metrics server: %v", err)
		return
	}
}

func RequestMW(name string, handler http.HandlerFunc) http.Handler {
	getExemplarFn := func(ctx context.Context) prometheus.Labels {
		if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsSampled() {
			return prometheus.Labels{"TraceID": spanCtx.TraceID().String()}
		}
		return nil
	}

	return promhttp.InstrumentHandlerInFlight(
		httpInFlightRequests,
		promhttp.InstrumentHandlerDuration(
			httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": name}),
			promhttp.InstrumentHandlerCounter(
				httpRequestsTotal.MustCurryWith(prometheus.Labels{"handler": name}),
				promhttp.InstrumentHandlerResponseSize(
					responseSize,
					handler,
					promhttp.WithExemplarFromContext(getExemplarFn),
				),
				promhttp.WithExemplarFromContext(getExemplarFn),
			),
			promhttp.WithExemplarFromContext(getExemplarFn),
		),
	)
}
