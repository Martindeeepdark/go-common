package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracer(t *testing.T) (*tracer, *tracetest.SpanRecorder) {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	tr := &tracer{
		provider: provider,
		tp:       provider.Tracer("test"),
	}
	saved := defaultTracer
	savedProvider := otel.GetTracerProvider()
	defaultTracer = tr
	otel.SetTracerProvider(provider)
	t.Cleanup(func() {
		defaultTracer = saved
		otel.SetTracerProvider(savedProvider)
		provider.Shutdown(context.Background())
	})
	return tr, rec
}

func TestHandler_CreatesSpan(t *testing.T) {
	tr, rec := setupTestTracer(t)

	var capturedTraceID string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTraceID = tr.TraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := Handler(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recw := httptest.NewRecorder()
	handler.ServeHTTP(recw, req)

	assert.Equal(t, http.StatusOK, recw.Code)
	assert.NotEmpty(t, capturedTraceID, "handler should create a trace span")

	spans := rec.Ended()
	require.Len(t, spans, 1)
	assert.Contains(t, spans[0].Name(), "http")
}

func TestHandler_PropagatesTraceID(t *testing.T) {
	_, rec := setupTestTracer(t)

	var ctxTraceID string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxTraceID = TraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := Handler(inner)

	req := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	recw := httptest.NewRecorder()
	handler.ServeHTTP(recw, req)

	assert.Equal(t, http.StatusOK, recw.Code)
	assert.NotEmpty(t, ctxTraceID)

	require.Len(t, rec.Ended(), 1)
	assert.Equal(t, ctxTraceID, rec.Ended()[0].SpanContext().TraceID().String())
}

func TestHandler_DifferentRequests_DifferentTraceIDs(t *testing.T) {
	tr, _ := setupTestTracer(t)

	var ids []string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids = append(ids, tr.TraceID(r.Context()))
		w.WriteHeader(http.StatusOK)
	})

	handler := Handler(inner)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		recw := httptest.NewRecorder()
		handler.ServeHTTP(recw, req)
	}

	assert.Len(t, ids, 3)
	assert.NotEqual(t, ids[0], ids[1])
	assert.NotEqual(t, ids[1], ids[2])
}

func TestHandler_StatusCodes(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"200 OK", http.StatusOK},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, rec := setupTestTracer(t)

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			})
			handler := Handler(inner)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			recw := httptest.NewRecorder()
			handler.ServeHTTP(recw, req)

			assert.Equal(t, tt.status, recw.Code)
			assert.Len(t, rec.Ended(), 1)
		})
	}
}

func TestGinMiddleware_CreatesSpan(t *testing.T) {
	tr, rec := setupTestTracer(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinMiddleware())

	var capturedTraceID string
	r.GET("/ping", func(c *gin.Context) {
		capturedTraceID = tr.TraceID(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	recw := httptest.NewRecorder()
	r.ServeHTTP(recw, req)

	assert.Equal(t, http.StatusOK, recw.Code)
	assert.NotEmpty(t, capturedTraceID, "gin middleware should create a trace span")
	assert.Len(t, rec.Ended(), 1)
}

func TestGinMiddleware_PropagatesTraceID(t *testing.T) {
	_, rec := setupTestTracer(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinMiddleware())

	var ctxTraceID string
	r.GET("/api/test", func(c *gin.Context) {
		ctxTraceID = TraceID(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	recw := httptest.NewRecorder()
	r.ServeHTTP(recw, req)

	assert.NotEmpty(t, ctxTraceID)
	require.Len(t, rec.Ended(), 1)
	assert.Equal(t, ctxTraceID, rec.Ended()[0].SpanContext().TraceID().String())
}

func TestGinMiddleware_MultipleRoutes(t *testing.T) {
	_, _ = setupTestTracer(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinMiddleware())

	var ids []string
	r.GET("/a", func(c *gin.Context) {
		ids = append(ids, TraceID(c.Request.Context()))
		c.Status(http.StatusOK)
	})
	r.POST("/b", func(c *gin.Context) {
		ids = append(ids, TraceID(c.Request.Context()))
		c.Status(http.StatusCreated)
	})

	for _, path := range []string{"/a", "/b"} {
		method := http.MethodGet
		if path == "/b" {
			method = http.MethodPost
		}
		req := httptest.NewRequest(method, path, nil)
		recw := httptest.NewRecorder()
		r.ServeHTTP(recw, req)
	}

	assert.Len(t, ids, 2)
	assert.NotEqual(t, ids[0], ids[1], "different requests should have different trace IDs")
}
