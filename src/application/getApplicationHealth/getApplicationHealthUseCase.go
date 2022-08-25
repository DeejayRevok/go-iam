package getApplicationHealth

import (
	"go-uaa/src/domain/healthcheck"
	"go-uaa/src/domain/internals"

	"go.uber.org/zap"
)

type GetApplicationHealthUseCase struct {
	healthChecker *healthcheck.HealthChecker
	logger        *zap.Logger
}

func (useCase *GetApplicationHealthUseCase) Execute(_ any) internals.UseCaseResponse {
	useCase.logger.Info("Starting checking if application is healthy")
	defer useCase.logger.Info("Finished checking if application is healthy")
	return internals.UseCaseResponse{
		Err: useCase.healthChecker.Check(),
	}
}

func (*GetApplicationHealthUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewGetApplicationHealthUseCase(healthChecker *healthcheck.HealthChecker, logger *zap.Logger) *GetApplicationHealthUseCase {
	useCase := GetApplicationHealthUseCase{
		healthChecker: healthChecker,
		logger:        logger,
	}
	return &useCase
}
