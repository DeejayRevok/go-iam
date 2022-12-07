package getThirdPartyAuthenticationUrl

import (
	"context"
	"errors"
	"go-iam/mocks"
	"go-iam/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	ThirdPartuAuthURLBuilderFactory *mocks.ThirdPartyAuthURLBuilderFactory
	UseCase                         *GetThirdPartyAuthenticationURLUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	urlBuilderFactoryMock := mocks.NewThirdPartyAuthURLBuilderFactory(t)
	return testCase{
		ThirdPartuAuthURLBuilderFactory: urlBuilderFactoryMock,
		UseCase:                         NewGetThirdPartyAuthenticationURLUseCase(urlBuilderFactoryMock, logger),
	}
}

func TestExecuteWrongRequest(t *testing.T) {
	testCase := setUp(t)
	request := "wrongRequest"
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.ThirdPartuAuthURLBuilderFactory.AssertNotCalled(t, "Create")
}

func TestExecuteCreateURLBuilderError(t *testing.T) {
	testCase := setUp(t)
	request := GetThirdPartyAuthenticationURLRequest{
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	createError := errors.New("Test create URL builder factory error")
	testCase.ThirdPartuAuthURLBuilderFactory.On("Create", mock.Anything).Return(nil, createError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != createError {
		t.Fatal("Expected use case to return URL builder factory error")
	}
	testCase.ThirdPartuAuthURLBuilderFactory.AssertCalled(t, "Create", request.AuthProvider)
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	request := GetThirdPartyAuthenticationURLRequest{
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testAuthURL := "testAuthURL"
	authURLBuilder := mocks.NewThirdPartyAuthURLBuilder(t)
	authURLBuilder.On("Build", mock.Anything).Return(testAuthURL)
	testCase.ThirdPartuAuthURLBuilderFactory.On("Create", mock.Anything).Return(authURLBuilder, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	responseAuthURL := response.Content.(string)
	if responseAuthURL != testAuthURL {
		t.Fatal("Expected use case to return same URL as the builder returned one")
	}
	testCase.ThirdPartuAuthURLBuilderFactory.AssertCalled(t, "Create", request.AuthProvider)
	authURLBuilder.AssertCalled(t, "Build", request.CallbackURL)
}
