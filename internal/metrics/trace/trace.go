package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type JaegerTracing struct {
	Exporter tracesdk.SpanExporter
	Provider *tracesdk.TracerProvider
	Tracer   trace.Tracer
}

func NewJaegerExporter(ctx context.Context) (tracesdk.SpanExporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}

func NewTraceProvider(exp tracesdk.SpanExporter, ServiceName string) (*tracesdk.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

func InitTracer(serviceName string) *JaegerTracing {
	exporter, err := NewJaegerExporter(context.Background())
	if err != nil {
		log.Fatalf("initialize tracer exporter: %v", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		log.Fatalf("initialize tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)

	tracer, err := tp.Tracer("main tracer"), nil

	if err != nil {
		log.Fatalf("error while init tracer: %v", err)
	}

	return &JaegerTracing{
		Exporter: exporter,
		Provider: tp,
		Tracer:   tracer,
	}
}
