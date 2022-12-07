package getApplicationHealth

import (
	"context"
	"go-iam/src/domain/healthcheck"
	"go-iam/src/domain/internals"
)

type GetApplicationHealthUseCase struct {
	healthChecker *healthcheck.HealthChecker
	logger        internals.Logger
}

func (useCase *GetApplicationHealthUseCase) Execute(ctx context.Context, _ any) internals.UseCaseResponse {
	useCase.logger.Info(ctx, "Starting checking if application is healthy")
	defer useCase.logger.Info(ctx, "Finished checking if application is healthy")
	return internals.UseCaseResponse{
		Err: useCase.healthChecker.Check(),
	}
}

func NewGetApplicationHealthUseCase(healthChecker *healthcheck.HealthChecker, logger internals.Logger) *GetApplicationHealthUseCase {
	useCase := GetApplicationHealthUseCase{
		healthChecker: healthChecker,
		logger:        logger,
	}
	return &useCase
}
