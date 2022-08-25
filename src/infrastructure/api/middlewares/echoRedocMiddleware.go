package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
)

func NewEchoRedocMiddleware(doc redoc.Redoc) echo.MiddlewareFunc {
	handle := doc.Handler()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			handle(ctx.Response(), ctx.Request())
			return next(ctx)
		}
	}
}
