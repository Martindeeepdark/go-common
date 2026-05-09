package otel

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type mockTransport struct {
	operation     string
	requestHeader transport.Header
}

func (m *mockTransport) Kind() transport.Kind            { return transport.KindGRPC }
func (m *mockTransport) Endpoint() string                { return "/test" }
func (m *mockTransport) Operation() string               { return m.operation }
func (m *mockTransport) RequestHeader() transport.Header { return m.requestHeader }
func (m *mockTransport) ReplyHeader() transport.Header   { return nil }

type mockHeader struct {
	data map[string][]string
}

func (h *mockHeader) Get(key string) string {
	if v, ok := h.data[key]; ok && len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h *mockHeader) Set(key, value string) {
	h.data[key] = []string{value}
}

func (h *mockHeader) Add(key, value string) {
	h.data[key] = append(h.data[key], value)
}

func (h *mockHeader) Values(key string) []string {
	return h.data[key]
}

func (h *mockHeader) Keys() []string {
	keys := make([]string, 0, len(h.data))
	for k := range h.data {
		keys = append(keys, k)
	}
	return keys
}

func setupKratosTest(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))

	savedTracer := defaultTracer
	savedProvider := otel.GetTracerProvider()
	defaultTracer = &tracer{provider: provider, tp: provider.Tracer("test")}
	otel.SetTracerProvider(provider)
	t.Cleanup(func() {
		defaultTracer = savedTracer
		otel.SetTracerProvider(savedProvider)
		provider.Shutdown(context.Background())
	})
	return rec
}

func TestKratosServer_CreatesSpan(t *testing.T) {
	rec := setupKratosTest(t)

	m := KratosServer()
	require.NotNil(t, m)

	header := &mockHeader{data: map[string][]string{}}
	ctx := transport.NewServerContext(context.Background(), &mockTransport{
		operation:     "test.operation",
		requestHeader: header,
	})

	handler := m(func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	reply, err := handler(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, "ok", reply)

	spans := rec.Ended()
	assert.NotEmpty(t, spans, "KratosServer middleware should create a span")
}

func TestKratosClient_CreatesSpan(t *testing.T) {
	rec := setupKratosTest(t)

	m := KratosClient()
	require.NotNil(t, m)

	header := &mockHeader{data: map[string][]string{}}
	ctx := transport.NewClientContext(context.Background(), &mockTransport{
		operation:     "client.operation",
		requestHeader: header,
	})

	handler := m(func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	reply, err := handler(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, "ok", reply)

	spans := rec.Ended()
	assert.NotEmpty(t, spans, "KratosClient middleware should create a span")
}

func TestKratosServer_NoTransport(t *testing.T) {
	_ = setupKratosTest(t)

	m := KratosServer()
	handler := m(func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	reply, err := handler(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, "ok", reply)
}
