package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4/v2"
)

func NewEchoAPMMiddleware() echo.MiddlewareFunc {
	return apmechov4.Middleware()
}
