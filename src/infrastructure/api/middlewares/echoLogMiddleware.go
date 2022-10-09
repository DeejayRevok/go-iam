package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type EchoLogMiddleware struct {
	logger *zap.Logger
}

func (middleware *EchoLogMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			request := c.Request()
			response := c.Response()
			if err != nil {
				middleware.logger.Error(err.Error())
				c.Error(err)
			}
			if request.URL.Path != "/metrics" {
				middleware.logger.Info(fmt.Sprintf("Executed HTTP %s request for %s with response status %d", request.Method, request.RequestURI, response.Status))
			}
			return nil
		}
	}
}

func NewEchoLogMiddleware(logger *zap.Logger) *EchoLogMiddleware {
	middleware := EchoLogMiddleware{
		logger: logger,
	}
	return &middleware
}
