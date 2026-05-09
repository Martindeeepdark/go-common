package otel

import (
	"context"
	"fmt"
	"os"
	"sync"

	"common/trace/defs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	defaultTracer *tracer
	once          sync.Once
)

var _ defs.Tracer = (*tracer)(nil)

type tracer struct {
	provider *sdktrace.TracerProvider
	tp       trace.Tracer
}

type Option func(*config)

type config struct {
	endpoint    string
	serviceName string
	insecure    bool
}

func WithEndpoint(addr string) Option {
	return func(c *config) { c.endpoint = addr }
}

func WithServiceName(name string) Option {
	return func(c *config) { c.serviceName = name }
}

func WithInsecure() Option {
	return func(c *config) { c.insecure = true }
}

func Init(opts ...Option) (func(context.Context) error, error) {
	var initErr error
	once.Do(func() {
		cfg := config{
			endpoint:    envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			serviceName: envOr("OTEL_SERVICE_NAME", "unknown"),
		}
		for _, opt := range opts {
			opt(&cfg)
		}

		expOpts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.endpoint)}
		if cfg.insecure {
			expOpts = append(expOpts, otlptracegrpc.WithInsecure())
		}

		exp, err := otlptracegrpc.New(context.Background(), expOpts...)
		if err != nil {
			initErr = fmt.Errorf("trace: create exporter: %w", err)
			return
		}

		provider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
		)
		otel.SetTracerProvider(provider)

		defaultTracer = &tracer{
			provider: provider,
			tp:       provider.Tracer(cfg.serviceName),
		}
	})
	if initErr != nil {
		return nil, initErr
	}

	return func(ctx context.Context) error {
		if defaultTracer != nil && defaultTracer.provider != nil {
			return defaultTracer.provider.Shutdown(ctx)
		}
		return nil
	}, nil
}

func (t *tracer) TraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasTraceID() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func (t *tracer) SpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasSpanID() {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

func TraceID(ctx context.Context) string {
	if defaultTracer == nil {
		return ""
	}
	return defaultTracer.TraceID(ctx)
}

func SpanID(ctx context.Context) string {
	if defaultTracer == nil {
		return ""
	}
	return defaultTracer.SpanID(ctx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
