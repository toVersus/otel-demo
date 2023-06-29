package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	meterName     = "github.com/toVersus/otel-demo"
	metricsPrefix = "oteldemo"
)

type CMetric int

const (
	ErrorCount CMetric = CMetric(iota)
)

type counterMetric struct {
	metricName    string
	metricDesc    string
	counter       api.Int64Counter
	createCounter *sync.Once
}

var counterMetricMap = map[CMetric]*counterMetric{
	ErrorCount: {metricsPrefix + "_errors", "Total number of errors made", nil, &errorCounterOnce},
}

var (
	errorCounterOnce sync.Once
)

func InitMetricProvider() {
	exporter, err := prometheus.New()
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
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
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

func GetCounter(name CMetric) (api.Int64Counter, error) {
	counterMetric, ok := counterMetricMap[name]
	if !ok {
		return nil, errors.New("counter not exists")
	}
	counterMetric.createCounter.Do(func() {
		meter := otel.Meter(meterName)
		counterMetric.counter, _ = meter.Int64Counter(counterMetric.metricName,
			api.WithDescription(counterMetric.metricDesc),
		)
	})
	if counterMetric.counter != nil {
		return counterMetric.counter, nil
	}
	return nil, errors.New("counter metric not ready")
}

func RecordCount(name CMetric, service string) error {
	ctx := context.Background()
	counter, err := GetCounter(name)
	if err != nil {
		return err
	}

	counter.Add(ctx, 1, api.WithAttributes(
		attribute.Key("service").String(service),
	))
	return nil
}
