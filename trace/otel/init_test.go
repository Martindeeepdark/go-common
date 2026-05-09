package otel

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	otlpcollectortrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// mockCollector is a minimal OTLP trace collector for testing Init.
type mockCollector struct {
	otlpcollectortrace.UnimplementedTraceServiceServer
}

func (m *mockCollector) Export(context.Context, *otlpcollectortrace.ExportTraceServiceRequest) (*otlpcollectortrace.ExportTraceServiceResponse, error) {
	return &otlpcollectortrace.ExportTraceServiceResponse{}, nil
}

func newMockCollector(t *testing.T) string {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	otlpcollectortrace.RegisterTraceServiceServer(s, &mockCollector{})
	go s.Serve(lis)
	t.Cleanup(s.Stop)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return conn.Target()
}

func TestInit_WithMockCollector(t *testing.T) {
	addr := newMockCollector(t)

	shutdown, err := Init(
		WithEndpoint(addr),
		WithServiceName("test-service"),
		WithInsecure(),
	)
	require.NoError(t, err)
	require.NotNil(t, shutdown)
	t.Cleanup(func() { shutdown(context.Background()) })

	// Verify defaultTracer is set
	require.NotNil(t, defaultTracer)
	assert.NotNil(t, defaultTracer.provider)
	assert.NotNil(t, defaultTracer.tp)
}

func TestInit_OnlyOnce(t *testing.T) {
	once = sync.Once{}
	defer func() { once = sync.Once{} }()

	addr := newMockCollector(t)

	tracerBefore := defaultTracer

	shutdown1, err1 := Init(WithEndpoint(addr), WithInsecure())
	require.NoError(t, err1)
	require.NotNil(t, shutdown1)

	tracerAfterFirst := defaultTracer

	// Second call should not change the tracer (sync.Once)
	_, err2 := Init(WithEndpoint("invalid:9999"))
	require.NoError(t, err2)

	assert.Equal(t, tracerAfterFirst, defaultTracer, "second Init should not change the tracer")
	assert.NotEqual(t, tracerBefore, tracerAfterFirst)

	shutdown1(context.Background())
}
