package getApplicationHealth

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/healthcheck"
	"go-uaa/src/infrastructure/logging"
	"testing"

	"go.elastic.co/apm/v2"
)

type testCase struct {
	SingleHealthChecker *mocks.SingleHealthChecker
	UseCase             *GetApplicationHealthUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	singleHealthCheckerMock := mocks.NewSingleHealthChecker(t)
	singleHealthCheckers := []healthcheck.SingleHealthChecker{singleHealthCheckerMock}
	healthChecker := healthcheck.NewHealthChecker(singleHealthCheckers)
	return testCase{
		SingleHealthChecker: singleHealthCheckerMock,
		UseCase:             NewGetApplicationHealthUseCase(healthChecker, logger),
	}
}

func TestExecuteHealthCheckError(t *testing.T) {
	testCase := setUp(t)
	checkError := errors.New("Test health check error")
	testCase.SingleHealthChecker.On("Check").Return(checkError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, struct{}{})

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != checkError {
		t.Fatal("Expected use case to return health check error")
	}
	testCase.SingleHealthChecker.AssertCalled(t, "Check")
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testCase.SingleHealthChecker.On("Check").Return(nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, struct{}{})

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.SingleHealthChecker.AssertCalled(t, "Check")
}
