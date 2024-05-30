package tracing

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	ServiceName string
	Tracer      trace.Tracer
}

var Tracer *TracingConfig

func NewTracingConfig(serviceName string) *TracingConfig {
	tracerProvider := initializeTracer(serviceName)
	Tracer = &TracingConfig{
		ServiceName: serviceName,
		Tracer:      tracerProvider.Tracer(serviceName),
	}
	return Tracer
}

// initializeTracer initializes a new tracer provider.
func initializeTracer(serviceName string) *sdktrace.TracerProvider {
	// Set up an OTLP exporter for Tempo
	expOTLP, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("tempo.tempo.svc.cluster.local:4317"))
	if err != nil {
		logrus.Fatalf("failed to create exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(expOTLP),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Register the trace provider as the global provider
	otel.SetTracerProvider(tp)

	return tp
}

// Middleware adds tracing to each request.
func (tc *TracingConfig) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tc.Tracer.Start(c.Context(), c.OriginalURL(), trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.url", c.OriginalURL()),
		))
		defer span.End()
		traceID := span.SpanContext().TraceID().String()
		c.Locals("trace_id", traceID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

// Lebih bersih: hanya menggunakan ctx untuk mengekstrak dan menetapkan span.
func (tc *TracingConfig) SetError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil && err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

type Span struct {
	span trace.Span
}

// StartTraceFromContext starts a new trace from the given context and returns the context and span.
func StartTraceFromContext(ctx context.Context, spanName string) (context.Context, *Span) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, spanName)
	return ctx, &Span{span}
}

func (s *Span) SetError(err error) {
	if s.span == nil || err == nil {
		return
	}
	s.span.RecordError(err, trace.WithStackTrace(true))
	s.span.SetStatus(codes.Error, err.Error())
	s.span.SetAttributes(attribute.String("error.message", err.Error()))
}
