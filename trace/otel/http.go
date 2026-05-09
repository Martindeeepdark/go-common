package otel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Handler(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http")
}

func GinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("gin")
}
