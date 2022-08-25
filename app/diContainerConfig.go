package app

import (
	"go-uaa/app/cli/commands"
	"go-uaa/src/application/authenticate"
	"go-uaa/src/application/createPermission"
	"go-uaa/src/application/createRole"
	"go-uaa/src/application/createUser"
	"go-uaa/src/application/getApplicationHealth"
	"go-uaa/src/application/getAuthenticatedUser"
	"go-uaa/src/application/getUser"
	"go-uaa/src/application/requestPasswordReset"
	"go-uaa/src/application/resetPassword"
	"go-uaa/src/application/sendPasswordResetToken"
	"go-uaa/src/domain/auth"
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/auth/authenticationStrategy"
	"go-uaa/src/domain/auth/refreshToken"
	"go-uaa/src/domain/events"
	"go-uaa/src/domain/hash"
	"go-uaa/src/domain/healthcheck"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/api/controllers"
	"go-uaa/src/infrastructure/api/middlewares"
	"go-uaa/src/infrastructure/database"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/email"
	"go-uaa/src/infrastructure/jwt"
	"go-uaa/src/infrastructure/messaging"
	"go-uaa/src/infrastructure/security"
	"go-uaa/src/infrastructure/transformers"

	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildDIContainer() dig.Container {
	container := dig.New()

	container.Provide(NewLogger)
	container.Provide(ConnectDatabase)
	container.Provide(ConnectToAMQPServer)
	container.Provide(LoadJWTSettings)
	container.Provide(BuildSMTPClient)

	container.Provide(database.NewPermissionDbRepository, dig.As(new(permission.PermissionRepository)))
	container.Provide(database.NewRoleDbRepository, dig.As(new(role.RoleRepository)))
	container.Provide(database.NewUserDbRepository, dig.As(new(user.UserRepository)))
	container.Provide(database.NewUserPasswordResetDbRepository, dig.As(new(user.UserPasswordResetRepository)))

	container.Provide(security.NewBcryptHasher, dig.As(new(hash.Hasher)))
	container.Provide(security.NewBcryptHashComparator, dig.As(new(hash.HashComparator)))
	container.Provide(security.NewMailCheckerEmailValidator, dig.As(new(user.EmailValidator)))
	container.Provide(jwt.NewJWTTokenGenerator)
	container.Provide(jwt.NewJWTClaimsToAccessTokenTransformer)
	container.Provide(jwt.NewJWTClaimsToRefreshTokenTransformer)
	container.Provide(jwt.NewJWTAccessTokenDeserializer, dig.As(new(accessToken.AccessTokenDeserializer)))
	container.Provide(jwt.NewJWTRefreshTokenDeserializer, dig.As(new(refreshToken.RefreshTokenDeserializer)))
	container.Provide(jwt.NewJWTRSAKeyToJWTKeyResponseTransformer)
	container.Provide(jwt.NewJWTKeySetBuilder)

	container.Provide(transformers.NewRoleToResponseTransformer)
	container.Provide(transformers.NewUserToResponseTransformer)
	container.Provide(transformers.NewPermissionToResponseTransformer)
	container.Provide(transformers.NewEventToAMQPMessageTransformer)
	container.Provide(transformers.NewAccessTokenToJWTClaimsTransformer)
	container.Provide(transformers.NewRefreshTokenToJWTClaimsTransformer)
	container.Provide(transformers.NewAuthenticationToResponseTransformer)
	container.Provide(transformers.NewAMQPDeliveryToMapTransformer)
	container.Provide(transformers.NewErrorToEchoErrorTransformer)

	container.Provide(accessToken.NewAccessTokenGenerator)
	container.Provide(refreshToken.NewRefreshTokenGenerator)
	container.Provide(authenticationStrategy.NewPasswordAuthenticationStrategy)
	container.Provide(authenticationStrategy.NewRefreshTokenAuthenticationStrategy)
	container.Provide(auth.NewAuthenticator)

	container.Provide(func(amqpConnection *amqp.Connection, logger *zap.Logger) *amqp.Channel {
		amqpChannel, err := amqpConnection.Channel()
		if err != nil {
			logger.Fatal("Error creating the AMQP channel")
			return nil
		}
		return amqpChannel
	})
	container.Provide(messaging.NewAMQPExchangeManager)
	container.Provide(messaging.NewAMQPExchangeEventPublisher, dig.As(new(events.EventPublisher)))
	container.Provide(messaging.NewAMQPQueueEventListenerFactory, dig.As(new(events.EventListenerFactory)))

	container.Provide(email.NewEmailPasswordResetTokenSender, dig.As(new(user.PasswordResetTokenSender)))

	container.Provide(internals.NewAuthorizedUseCaseExecutor)
	container.Provide(createUser.NewCreateUserUseCase)
	container.Provide(getUser.NewGetUserUseCase)
	container.Provide(createPermission.NewCreatePermissionUseCase)
	container.Provide(createRole.NewCreateRoleUseCase)
	container.Provide(authenticate.NewAuthenticationUseCase)
	container.Provide(getAuthenticatedUser.NewGetAuthenticatedUserUseCase)
	container.Provide(requestPasswordReset.NewRequestPasswordResetUseCase)
	container.Provide(sendPasswordResetToken.NewSendPasswordResetTokenUseCase)
	container.Provide(sendPasswordResetToken.NewUserPasswordResetRequestedConsumer)
	container.Provide(resetPassword.NewResetPasswordUseCase)

	container.Provide(dto.NewEchoDTOSerializer)
	container.Provide(dto.NewEchoDTODeserializer)

	container.Provide(dto.NewDTOValidator)
	container.Provide(middlewares.NewEchoLogMiddleware)
	container.Provide(NewRedocConfiguration)

	addHealthCheckDependencies(container)

	container.Provide(api.NewHTTPAccessTokenFinder)
	container.Provide(controllers.NewCreateUserController)
	container.Provide(controllers.NewGetUserController)
	container.Provide(controllers.NewCreatePermissionController)
	container.Provide(controllers.NewCreateRoleController)
	container.Provide(controllers.NewAuthenticateController)
	container.Provide(controllers.NewGetAuthenticatedUserController)
	container.Provide(controllers.NewGetJWTKeySetController)
	container.Provide(controllers.NewRequestResetPasswordController)
	container.Provide(controllers.NewResetPasswordController)

	container.Provide(commands.NewBoostrapPermissionsCLI)
	container.Provide(commands.NewCreateSuperuserCLI)

	return *container
}

func addHealthCheckDependencies(diContainer *dig.Container) {
	type healthCheckersAggregator struct {
		dig.Out
		HealthChecker healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	diContainer.Provide(func(db *gorm.DB) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: database.NewDatabaseHealthChecker(db),
		}
	})
	diContainer.Provide(func(amqpConnection *amqp.Connection) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: messaging.NewAMQPHealthChecker(amqpConnection),
		}
	})

	type healthCheckersGroup struct {
		dig.In
		HealthCheckers []healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	diContainer.Provide(func(checkersGroup healthCheckersGroup) *healthcheck.HealthChecker {
		return healthcheck.NewHealthChecker(checkersGroup.HealthCheckers)
	})

	diContainer.Provide(getApplicationHealth.NewGetApplicationHealthUseCase)
	diContainer.Provide(controllers.NewGetStatusController)
}
