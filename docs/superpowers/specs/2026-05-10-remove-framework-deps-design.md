# Remove Framework Dependencies from Common Package

**Date:** 2026-05-10
**Status:** Approved

## Problem

The `common` package directly depends on Gin v1.12.0, Kratos v2.9.2, and gRPC. This forces all consumers to pull in these frameworks and their transitive dependencies, even when they don't use them. It also locks framework versions across all services.

## Decision

Remove all framework-specific adapter code from `common`. Keep only framework-agnostic utilities.

## Files to Delete

- `trace/otel/http.go` — Gin middleware + net/http handler
- `trace/otel/http_test.go`
- `trace/otel/grpc.go` — gRPC server/client options
- `trace/otel/grpc_test.go`
- `trace/otel/kratos.go` — Kratos tracing middleware
- `trace/otel/kratos_test.go`
- `trace/otel/init_test.go` — depends on gRPC bufconn mock

## Files to Modify

- `trace/defs/tracer.go` — Remove `HTTPMiddleware` and `GRPCInterceptor` interfaces; keep only `Tracer`

## Files Unchanged

- `trace/otel/tracer.go` — Core Init/TraceID/SpanID (only depends on otel SDK)
- `trace/otel/tracer_test.go` — Core unit tests (no framework deps)

## go.mod Changes

### Remove direct dependencies

- `github.com/gin-gonic/gin`
- `github.com/go-kratos/kratos/v2`
- `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`
- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`
- `google.golang.org/grpc`

### Keep

- All `go.opentelemetry.io/otel*` packages (core trace SDK)
- All other packages: redis, zap, sonyflake, testify, yaml, etc.

## Post-Cleanup

Services that need framework-specific tracing should import the official adapter libraries directly:
- `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`
- `github.com/go-kratos/kratos/v2/middleware/tracing`
