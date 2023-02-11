package app

import (
	"fmt"
	"go-iam/app/cli/commands"
	"go-iam/src/application/authenticate"
	"go-iam/src/application/createUser"
	"go-iam/src/application/getApplicationHealth"
	"go-iam/src/application/getAuthenticatedUser"
	"go-iam/src/application/getThirdPartyAuthenticationUrl"
	"go-iam/src/application/requestPasswordReset"
	"go-iam/src/application/resetPassword"
	"go-iam/src/application/sendPasswordResetToken"
	"go-iam/src/domain/auth"
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/auth/authenticationStrategy"
	"go-iam/src/domain/auth/refreshToken"
	"go-iam/src/domain/auth/thirdParty"
	"go-iam/src/domain/events"
	"go-iam/src/domain/hash"
	"go-iam/src/domain/healthcheck"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/api"
	"go-iam/src/infrastructure/api/controllers"
	"go-iam/src/infrastructure/api/middlewares"
	"go-iam/src/infrastructure/database"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/email"
	"go-iam/src/infrastructure/graph/resolvers"
	"go-iam/src/infrastructure/jwt"
	"go-iam/src/infrastructure/logging"
	"go-iam/src/infrastructure/messaging"
	"go-iam/src/infrastructure/oauth2"
	"go-iam/src/infrastructure/security"
	"go-iam/src/infrastructure/transformers"

	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildDIContainer() dig.Container {
	container := dig.New()
	if err := container.Provide(NewAPMTracer); err != nil {
		panic(fmt.Sprintf("Error providing APM tracer to the dependency injection container: %s", err.Error()))
	}
	if err := container.Provide(logging.NewZapLogger); err != nil {
		panic(fmt.Sprintf("Error providing zap logger to the dependency injection container: %s", err.Error()))
	}
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Provide(logging.NewZapTracedLogger, dig.As(new(internals.Logger))), logger)
		handleError(container.Provide(logging.NewZapGormTracedLogger), logger)
		handleError(container.Provide(ConnectDatabase), logger)
		handleError(container.Provide(ConnectToAMQPServer), logger)
		handleError(container.Provide(LoadJWTSettings), logger)
		handleError(container.Provide(BuildSMTPClient), logger)
		handleError(container.Provide(BuildOauth2GoogleAuthURLBuilder), logger)
		handleError(container.Provide(BuildOauth2GoogleTokensFetcher), logger)
		handleError(container.Provide(BuildThirdPartyAuthStateChecker), logger)
		handleError(container.Provide(BuildHTTPThirdPartyCallbackUrlBuilder), logger)

		handleError(container.Provide(database.NewUserDbRepository, dig.As(new(user.UserRepository))), logger)
		handleError(container.Provide(database.NewUserPasswordResetDbRepository, dig.As(new(user.UserPasswordResetRepository))), logger)

		handleError(container.Provide(security.NewBcryptHasher, dig.As(new(hash.Hasher))), logger)
		handleError(container.Provide(security.NewBcryptHashComparator, dig.As(new(hash.HashComparator))), logger)
		handleError(container.Provide(security.NewMailCheckerEmailValidator, dig.As(new(user.EmailValidator))), logger)

		handleError(container.Provide(jwt.NewJWTTokenGenerator), logger)
		handleError(container.Provide(jwt.NewJWTClaimsToAccessTokenTransformer), logger)
		handleError(container.Provide(jwt.NewJWTClaimsToRefreshTokenTransformer), logger)
		handleError(container.Provide(jwt.NewJWTAccessTokenDeserializer, dig.As(new(accessToken.AccessTokenDeserializer))), logger)
		handleError(container.Provide(jwt.NewJWTRefreshTokenDeserializer, dig.As(new(refreshToken.RefreshTokenDeserializer))), logger)
		handleError(container.Provide(jwt.NewJWTRSAKeyToJWTKeyResponseTransformer), logger)
		handleError(container.Provide(jwt.NewJWTKeySetBuilder), logger)
		handleError(container.Provide(jwt.NewJWTThirdPartyTokensToEmailTransformer, dig.As(new(thirdParty.ThirdPartyTokensToEmailTransformer))), logger)

		handleError(container.Provide(transformers.NewUserToResponseTransformer), logger)
		handleError(container.Provide(transformers.NewEventToAMQPMessageTransformer), logger)
		handleError(container.Provide(transformers.NewAccessTokenToJWTClaimsTransformer), logger)
		handleError(container.Provide(transformers.NewRefreshTokenToJWTClaimsTransformer), logger)
		handleError(container.Provide(transformers.NewAuthenticationToResponseTransformer), logger)
		handleError(container.Provide(transformers.NewAMQPDeliveryToMapTransformer), logger)
		handleError(container.Provide(transformers.NewErrorToEchoErrorTransformer), logger)
		handleError(container.Provide(transformers.NewOauth2TokenToThirdPartyTokensTransformer), logger)
		handleError(container.Provide(accessToken.NewAccessTokenGenerator), logger)
		handleError(container.Provide(refreshToken.NewRefreshTokenGenerator), logger)
		handleError(container.Provide(authenticationStrategy.NewPasswordAuthenticationStrategy), logger)
		handleError(container.Provide(authenticationStrategy.NewRefreshTokenAuthenticationStrategy), logger)
		handleError(container.Provide(authenticationStrategy.NewThirdPartyAuthenticationStrategy), logger)
		handleError(container.Provide(auth.NewAuthenticator), logger)

		handleError(container.Provide(oauth2.NewOauth2ThirdPartyAuthURLBuilderFactory, dig.As(new(thirdParty.ThirdPartyAuthURLBuilderFactory))), logger)
		handleError(container.Provide(oauth2.NewOauth2ThirdPartyTokensFetcherFactory, dig.As(new(thirdParty.ThirdPartyTokensFetcherFactory))), logger)

		handleError(container.Provide(func(amqpConnection *amqp.Connection, logger *zap.Logger) *amqp.Channel {
			amqpChannel, err := amqpConnection.Channel()
			if err != nil {
				logger.Fatal("Error creating the AMQP channel")
				return nil
			}
			return amqpChannel
		}), logger)
		handleError(container.Provide(messaging.NewAMQPExchangeManager), logger)
		handleError(container.Provide(messaging.NewAMQPExchangeEventPublisher, dig.As(new(events.EventPublisher))), logger)
		handleError(container.Provide(messaging.NewAMQPQueueEventListenerFactory, dig.As(new(events.EventListenerFactory))), logger)

		handleError(container.Provide(email.NewEmailPasswordResetTokenSender, dig.As(new(user.PasswordResetTokenSender))), logger)

		handleError(container.Provide(internals.NewAuthorizedUseCaseExecutor), logger)
		handleError(container.Provide(createUser.NewCreateUserUseCase), logger)
		handleError(container.Provide(authenticate.NewAuthenticationUseCase), logger)
		handleError(container.Provide(getAuthenticatedUser.NewGetAuthenticatedUserUseCase), logger)
		handleError(container.Provide(requestPasswordReset.NewRequestPasswordResetUseCase), logger)
		handleError(container.Provide(sendPasswordResetToken.NewSendPasswordResetTokenUseCase), logger)
		handleError(container.Provide(sendPasswordResetToken.NewUserPasswordResetRequestedConsumer), logger)
		handleError(container.Provide(resetPassword.NewResetPasswordUseCase), logger)
		handleError(container.Provide(getThirdPartyAuthenticationUrl.NewGetThirdPartyAuthenticationURLUseCase), logger)

		handleError(container.Provide(dto.NewEchoDTOSerializer), logger)
		handleError(container.Provide(dto.NewEchoDTODeserializer), logger)

		handleError(container.Provide(dto.NewDTOValidator), logger)
		handleError(container.Provide(middlewares.NewEchoLogMiddleware), logger)
		handleError(container.Provide(NewRedocConfiguration), logger)

		addHealthCheckDependencies(container, logger)

		handleError(container.Provide(api.NewHTTPAccessTokenFinder), logger)
		handleError(container.Provide(controllers.NewCreateUserController), logger)
		handleError(container.Provide(controllers.NewAuthenticateController), logger)
		handleError(container.Provide(controllers.NewGetAuthenticatedUserController), logger)
		handleError(container.Provide(controllers.NewGetJWTKeySetController), logger)
		handleError(container.Provide(controllers.NewRequestResetPasswordController), logger)
		handleError(container.Provide(controllers.NewResetPasswordController), logger)
		handleError(container.Provide(controllers.NewGetThirdPartyAuthenticationController), logger)
		handleError(container.Provide(controllers.NewThirdPartyAuthenticationCallbackController), logger)

		handleError(container.Provide(resolvers.NewMeResolver), logger)
		handleError(container.Provide(resolvers.NewCreateUserResolver), logger)
		handleError(container.Provide(resolvers.NewRootResolver), logger)

		handleError(container.Provide(commands.NewCreateSuperuserCLI), logger)
	}); err != nil {
		panic(fmt.Sprintf("Error adding dependencies to the container: %s", err.Error()))
	}

	return *container
}

func addHealthCheckDependencies(diContainer *dig.Container, logger *zap.Logger) {
	type healthCheckersAggregator struct {
		dig.Out
		HealthChecker healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	handleError(diContainer.Provide(func(db *gorm.DB) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: database.NewDatabaseHealthChecker(db),
		}
	}), logger)
	handleError(diContainer.Provide(func(amqpConnection *amqp.Connection) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: messaging.NewAMQPHealthChecker(amqpConnection),
		}
	}), logger)

	type healthCheckersGroup struct {
		dig.In
		HealthCheckers []healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	handleError(diContainer.Provide(func(checkersGroup healthCheckersGroup) *healthcheck.HealthChecker {
		return healthcheck.NewHealthChecker(checkersGroup.HealthCheckers)
	}), logger)

	handleError(diContainer.Provide(getApplicationHealth.NewGetApplicationHealthUseCase), logger)
	handleError(diContainer.Provide(controllers.NewGetStatusController), logger)
}
