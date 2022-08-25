package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
)

func NewEchoCorsMiddleware() echo.MiddlewareFunc {
	handlerCORS := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowedHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "api_key"},
	}).HandlerFunc
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			handlerCORS(ctx.Response(), ctx.Request())
			return next(ctx)
		}
	}
}
