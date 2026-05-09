package defs

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type Tracer interface {
	TraceID(ctx context.Context) string
	SpanID(ctx context.Context) string
}

type HTTPMiddleware interface {
	Handler(next http.Handler) http.Handler
}

type GRPCInterceptor interface {
	ServerOptions() []grpc.ServerOption
	ClientOptions() []grpc.DialOption
}
