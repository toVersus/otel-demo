package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(ctx context.Context, otlpAddr, serviceName, version string) (func(), error) {
	exporter, err := newExporter(ctx, otlpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exporter %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(newResource(serviceName, version)),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}, nil
}

func newExporter(ctx context.Context, otlpAddr string) (trace.SpanExporter, error) {
	var err error
	var exporter trace.SpanExporter

	if len(otlpAddr) == 0 {
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
			stdouttrace.WithPrettyPrint(),
		)
	} else {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(otlpAddr),
		)
	}
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func newResource(serviceName, version string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
