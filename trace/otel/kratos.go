package otel

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

func KratosServer() middleware.Middleware {
	return tracing.Server()
}

func KratosClient() middleware.Middleware {
	return tracing.Client()
}
