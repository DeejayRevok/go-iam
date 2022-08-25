package app

import (
	"go-uaa/src/infrastructure/api/controllers"
	"go-uaa/src/infrastructure/api/middlewares"
	"go-uaa/src/infrastructure/dto"

	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildHTTPServer(container *dig.Container) *echo.Echo {
	server := echo.New()
	server.Use(middlewares.NewEchoCorsMiddleware())

	container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(validator *dto.DTOValidator) {
			server.Validator = validator
		}), logger)
		handleError(container.Invoke(func(redoc *redoc.Redoc) {
			server.Use(middlewares.NewEchoRedocMiddleware(*redoc))
		}), logger)
		handleError(container.Invoke(func(middleware *middlewares.EchoLogMiddleware) {
			server.Use(middleware.Middleware())
		}), logger)

		handleError(container.Invoke(func(controller *controllers.CreateUserController) {
			server.POST("/users", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetUserController) {
			server.GET("/users/:id", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetAuthenticatedUserController) {
			server.GET("/users/me", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.CreateRoleController) {
			server.POST("/roles", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.CreatePermissionController) {
			server.POST("/permissions", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.AuthenticateController) {
			server.POST("/token", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetJWTKeySetController) {
			server.GET("/jwks", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetStatusController) {
			server.GET("/status", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.RequestResetPasswordController) {
			server.POST("/users/password/reset", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.ResetPasswordController) {
			server.PUT("/users/password", controller.Handle)
		}), logger)
	})

	return server
}
