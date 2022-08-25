package getApplicationHealth

import (
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/healthcheck"
	"testing"

	"go.uber.org/zap"
)

type testCase struct {
	SingleHealthChecker *mocks.SingleHealthChecker
	UseCase             *GetApplicationHealthUseCase
}

func setUp(t *testing.T) testCase {
	logger, _ := zap.NewDevelopment()
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

	response := testCase.UseCase.Execute(struct{}{})

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

	response := testCase.UseCase.Execute(struct{}{})

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.SingleHealthChecker.AssertCalled(t, "Check")
}
