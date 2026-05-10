package defs

import "context"

type Tracer interface {
	TraceID(ctx context.Context) string
	SpanID(ctx context.Context) string
}
