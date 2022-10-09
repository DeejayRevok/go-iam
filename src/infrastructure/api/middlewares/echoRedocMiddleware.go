package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
)

const RedocDocsPath = "/docs"
const RedocSpecPath = "/openapi.yaml"

func NewEchoRedocMiddleware(doc redoc.Redoc) echo.MiddlewareFunc {
	handle := doc.Handler()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()
			if request.URL.Path == RedocDocsPath || request.URL.Path == RedocSpecPath {
				handle(ctx.Response(), request)
				return nil
			}
			return next(ctx)
		}
	}
}
