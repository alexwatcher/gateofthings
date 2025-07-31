package telemetry

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace/noop"
)

func MustCreateResource(serviceName string, version string, env string) *resource.Resource {
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(version),
		attribute.String("env", env),
	))
	if err != nil {
		panic(err)
	}
	return res
}

func MustInitTracer(ctx context.Context, res *resource.Resource, traceEndpoint string) {
	if traceEndpoint == "" {
		slog.Warn("no trace endpoint provided")
		otel.SetTextMapPropagator(propagation.TraceContext{})
		otel.SetTracerProvider(noop.NewTracerProvider())
		return
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(traceEndpoint),
		otlptracegrpc.WithCompressor("gzip"),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter),
	)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
}

func MustInitMeter(ctx context.Context, res *resource.Resource, metricsEndpoint string) {
	if metricsEndpoint == "" {
		slog.Warn("no metrics endpoint provided")
		reader := metric.NewManualReader()
		provider := metric.NewMeterProvider(metric.WithReader(reader))
		otel.SetMeterProvider(provider)
		return
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(metricsEndpoint),
		otlpmetricgrpc.WithCompressor("gzip"),
	)
	if err != nil {
		panic(err)
	}
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(time.Second))),
	)
	otel.SetMeterProvider(mp)
}

func MustInitLogger(ctx context.Context, res *resource.Resource, logsEndpoint string) {
	if logsEndpoint == "" {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Warn("no logs endpoint provided")
		return
	}

	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(logsEndpoint),
		otlploggrpc.WithCompressor("gzip"),
	)
	if err != nil {
		panic(err)
	}
	lp := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
	)
	global.SetLoggerProvider(lp)

	logger := otelslog.NewLogger("sys")
	slog.SetDefault(logger)
}
