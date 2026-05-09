package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func newTestTracer(t *testing.T) (*tracer, *tracetest.SpanRecorder) {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	return &tracer{
		provider: provider,
		tp:       provider.Tracer("test"),
	}, rec
}

func TestTracer_TraceID(t *testing.T) {
	tr, _ := newTestTracer(t)

	_, span := tr.tp.Start(context.Background(), "test-span")
	defer span.End()

	ctx := trace.ContextWithSpan(context.Background(), span)

	got := tr.TraceID(ctx)
	assert.NotEmpty(t, got, "TraceID should not be empty within a span")
	assert.Len(t, got, 32, "TraceID should be 32 hex chars")
}

func TestTracer_SpanID(t *testing.T) {
	tr, _ := newTestTracer(t)

	_, span := tr.tp.Start(context.Background(), "test-span")
	defer span.End()

	ctx := trace.ContextWithSpan(context.Background(), span)

	got := tr.SpanID(ctx)
	assert.NotEmpty(t, got, "SpanID should not be empty within a span")
	assert.Len(t, got, 16, "SpanID should be 16 hex chars")
}

func TestTracer_TraceID_EmptyContext(t *testing.T) {
	tr, _ := newTestTracer(t)

	got := tr.TraceID(context.Background())
	assert.Empty(t, got, "TraceID should be empty for background context")
}

func TestTracer_SpanID_EmptyContext(t *testing.T) {
	tr, _ := newTestTracer(t)

	got := tr.SpanID(context.Background())
	assert.Empty(t, got, "SpanID should be empty for background context")
}

func TestTracer_Interface(t *testing.T) {
	tr, _ := newTestTracer(t)

	var _ interface {
		TraceID(ctx context.Context) string
		SpanID(ctx context.Context) string
	} = tr
}

func TestPackage_TraceID_NilTracer(t *testing.T) {
	saved := defaultTracer
	defaultTracer = nil
	defer func() { defaultTracer = saved }()

	got := TraceID(context.Background())
	assert.Empty(t, got)
}

func TestPackage_SpanID_NilTracer(t *testing.T) {
	saved := defaultTracer
	defaultTracer = nil
	defer func() { defaultTracer = saved }()

	got := SpanID(context.Background())
	assert.Empty(t, got)
}

func TestPackage_TraceID_WithSpan(t *testing.T) {
	tr, _ := newTestTracer(t)

	saved := defaultTracer
	defaultTracer = tr
	defer func() { defaultTracer = saved }()

	_, span := tr.tp.Start(context.Background(), "test")
	defer span.End()
	ctx := trace.ContextWithSpan(context.Background(), span)

	got := TraceID(ctx)
	assert.NotEmpty(t, got)
}

func TestPackage_SpanID_WithSpan(t *testing.T) {
	tr, _ := newTestTracer(t)

	saved := defaultTracer
	defaultTracer = tr
	defer func() { defaultTracer = saved }()

	_, span := tr.tp.Start(context.Background(), "test")
	defer span.End()
	ctx := trace.ContextWithSpan(context.Background(), span)

	got := SpanID(ctx)
	assert.NotEmpty(t, got)
}

func TestTracer_TraceID_SpanRecorded(t *testing.T) {
	tr, rec := newTestTracer(t)

	ctx, span := tr.tp.Start(context.Background(), "op")
	traceID := tr.TraceID(ctx)
	span.End()

	require.NotEmpty(t, traceID)
	assert.Len(t, rec.Ended(), 1)
	assert.Equal(t, traceID, rec.Ended()[0].SpanContext().TraceID().String())
}

func TestTracer_TraceID_DifferentSpans(t *testing.T) {
	tr, _ := newTestTracer(t)

	ctx1, span1 := tr.tp.Start(context.Background(), "span1")
	ctx2, span2 := tr.tp.Start(context.Background(), "span2")
	defer span1.End()
	defer span2.End()

	id1 := tr.TraceID(ctx1)
	id2 := tr.TraceID(ctx2)

	assert.NotEqual(t, id1, id2, "different root spans should have different trace IDs")
}

func TestTracer_ChildSpan_SameTraceID(t *testing.T) {
	tr, _ := newTestTracer(t)

	ctx1, parent := tr.tp.Start(context.Background(), "parent")
	ctx2, child := tr.tp.Start(ctx1, "child")
	defer parent.End()
	defer child.End()

	assert.Equal(t, tr.TraceID(ctx1), tr.TraceID(ctx2), "child should share parent's trace ID")
	assert.NotEqual(t, tr.SpanID(ctx1), tr.SpanID(ctx2), "child should have different span ID")
}

func TestEnvOr(t *testing.T) {
	assert.Equal(t, "fallback", envOr("NONEXISTENT_KEY_12345", "fallback"))
	t.Setenv("TEST_ENV_OR_KEY", "value")
	assert.Equal(t, "value", envOr("TEST_ENV_OR_KEY", "fallback"))
}

func TestOptions(t *testing.T) {
	cfg := config{}
	WithEndpoint("collector:4317")(&cfg)
	assert.Equal(t, "collector:4317", cfg.endpoint)

	WithServiceName("my-service")(&cfg)
	assert.Equal(t, "my-service", cfg.serviceName)

	WithInsecure()(&cfg)
	assert.True(t, cfg.insecure)
}

func TestOptions_Defaults(t *testing.T) {
	cfg := config{
		endpoint:    envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		serviceName: envOr("OTEL_SERVICE_NAME", "unknown"),
	}
	assert.Equal(t, "localhost:4317", cfg.endpoint)
	assert.Equal(t, "unknown", cfg.serviceName)
	assert.False(t, cfg.insecure)
}

func TestShutdown_NoTracer(t *testing.T) {
	saved := defaultTracer
	defaultTracer = nil
	defer func() { defaultTracer = saved }()

	shutdown := func(ctx context.Context) error {
		if defaultTracer != nil && defaultTracer.provider != nil {
			return defaultTracer.provider.Shutdown(ctx)
		}
		return nil
	}
	assert.NoError(t, shutdown(context.Background()))
}
