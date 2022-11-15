package app

import (
	"go-uaa/src/infrastructure/api/controllers"
	"go-uaa/src/infrastructure/api/middlewares"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/graph/resolvers"

	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildHTTPServer(container *dig.Container) *echo.Echo {
	server := echo.New()
	server.Use(middlewares.NewEchoCorsMiddleware())
	server.Use(middlewares.NewEchoAPMMiddleware())

	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(validator *dto.DTOValidator) {
			server.Validator = validator
		}), logger)
		handleError(container.Invoke(func(middleware *middlewares.EchoLogMiddleware) {
			server.Use(middleware.Middleware())
		}), logger)
		handleError(container.Invoke(func(redoc *redoc.Redoc) {
			server.Use(middlewares.NewEchoRedocMiddleware(*redoc))
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
			server.POST("/auth/token", controller.Handle)
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
		handleError(container.Invoke(func(resolver *resolvers.RootResolver) {
			handler := BuildGraphQLHTTPHandler(resolver, logger)
			server.POST("/graphql", handler)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetThirdPartyAuthenticationController) {
			server.GET("/auth/third-party/:provider", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.ThirdPartyAuthenticationCallbackController) {
			server.GET("/auth/:provider/callback", controller.Handle)
		}), logger)
	}); err != nil {
		panic("Error adding HTTP API components to the dependency injection container")
	}

	return server
}
