package middlewares

import (
	"go-uaa/src/infrastructure/tracing"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
)

func NewEchoJaegerMiddleware(config *tracing.JaegerTracerConfig) echo.MiddlewareFunc {
	return jaegertracing.TraceWithConfig(jaegertracing.TraceConfig{
		Tracer:  *config.Tracer,
		Skipper: nil,
	})
}
