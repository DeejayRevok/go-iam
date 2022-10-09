package thirdPartyAuthentication

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/domain/session"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/logging"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	TokensFetcherFactory     *mocks.ThirdPartyTokensFetcherFactory
	TokensToEmailTransformer *mocks.ThirdPartyTokensToEmailTransformer
	UserRepo                 *mocks.UserRepository
	InternalState            string
	UseCase                  *ThirdPartyAuthenticationUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	tokensFetcherFactoryMock := mocks.NewThirdPartyTokensFetcherFactory(t)
	tokensToEmailTransformerMock := mocks.NewThirdPartyTokensToEmailTransformer(t)
	userRepoMock := mocks.NewUserRepository(t)
	testInternalState := "testInternalState"
	return testCase{
		TokensFetcherFactory:     tokensFetcherFactoryMock,
		TokensToEmailTransformer: tokensToEmailTransformerMock,
		UserRepo:                 userRepoMock,
		InternalState:            testInternalState,
		UseCase:                  NewThirdPartyAuthenticationUseCase(thirdParty.NewThirdPartyAuthStateChecker(testInternalState), tokensFetcherFactoryMock, tokensToEmailTransformerMock, userRepoMock, logger),
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
	testCase.TokensFetcherFactory.AssertNotCalled(t, "Create")
	testCase.TokensToEmailTransformer.AssertNotCalled(t, "Transform")
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteAuthStateCheckError(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        "wrongState",
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err.Error() != "auth state is not valid" {
		t.Fatal("Expected use case to return auth state not valid error")
	}
	testCase.TokensFetcherFactory.AssertNotCalled(t, "Create")
	testCase.TokensToEmailTransformer.AssertNotCalled(t, "Transform")
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteTokensFetcherCreationError(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testCreationError := errors.New("Test tokens fetcher creation error")
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(nil, testCreationError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testCreationError {
		t.Fatal("Expected use to return same error as the tokens fetcher factory")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	testCase.TokensToEmailTransformer.AssertNotCalled(t, "Transform")
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteTokensFetchError(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testFetchError := errors.New("Test tokens fetch error")
	tokensFetcherMock := mocks.NewThirdPartyTokensFetcher(t)
	tokensFetcherMock.On("Fetch", mock.Anything, mock.Anything).Return(nil, testFetchError)
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(tokensFetcherMock, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testFetchError {
		t.Fatal("Expected use to return same error as the tokens fetcher")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.Code, request.CallbackURL)
	testCase.TokensToEmailTransformer.AssertNotCalled(t, "Transform")
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteTokensTransformError(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	transformError := errors.New("Test tokens to email transformation error")
	testTokens := &thirdParty.ThirdPartyTokens{
		AccessToken:  "testAccessToken",
		RefreshToken: "testRefreshToken",
		IDToken:      "testIDToken",
	}
	tokensFetcherMock := mocks.NewThirdPartyTokensFetcher(t)
	tokensFetcherMock.On("Fetch", mock.Anything, mock.Anything).Return(testTokens, nil)
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(tokensFetcherMock, nil)
	testCase.TokensToEmailTransformer.On("Transform", mock.Anything).Return("", transformError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != transformError {
		t.Fatal("Expected use to return same error as the tokens to email transformer")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.Code, request.CallbackURL)
	testCase.TokensToEmailTransformer.AssertCalled(t, "Transform", testTokens)
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteFindByEmailError(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testTokens := &thirdParty.ThirdPartyTokens{
		AccessToken:  "testAccessToken",
		RefreshToken: "testRefreshToken",
		IDToken:      "testIDToken",
	}
	tokensFetcherMock := mocks.NewThirdPartyTokensFetcher(t)
	tokensFetcherMock.On("Fetch", mock.Anything, mock.Anything).Return(testTokens, nil)
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(tokensFetcherMock, nil)
	testEmail := "testEmail"
	testCase.TokensToEmailTransformer.On("Transform", mock.Anything).Return(testEmail, nil)
	testFindError := errors.New("Test find error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, testFindError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testFindError {
		t.Fatal("Expected use to return same error as the user repo find by email")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.Code, request.CallbackURL)
	testCase.TokensToEmailTransformer.AssertCalled(t, "Transform", testTokens)
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, testEmail)
}

func TestExecutUserNotFound(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testTokens := &thirdParty.ThirdPartyTokens{
		AccessToken:  "testAccessToken",
		RefreshToken: "testRefreshToken",
		IDToken:      "testIDToken",
	}
	tokensFetcherMock := mocks.NewThirdPartyTokensFetcher(t)
	tokensFetcherMock.On("Fetch", mock.Anything, mock.Anything).Return(testTokens, nil)
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(tokensFetcherMock, nil)
	testEmail := "testEmail"
	testCase.TokensToEmailTransformer.On("Transform", mock.Anything).Return(testEmail, nil)
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err.Error() != "user with email testEmail not registered" {
		t.Fatal("Expected use to return user not found error")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.Code, request.CallbackURL)
	testCase.TokensToEmailTransformer.AssertCalled(t, "Transform", testTokens)
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, testEmail)
}

func TestExecutSuccess(t *testing.T) {
	testCase := setUp(t)
	request := ThirdPartyAuthenticationRequest{
		State:        testCase.InternalState,
		Code:         "testCode",
		AuthProvider: "testAuthProvider",
		CallbackURL:  "testCallbackURL",
	}
	ctx := context.Background()
	testTokens := &thirdParty.ThirdPartyTokens{
		AccessToken:  "testAccessToken",
		RefreshToken: "testRefreshToken",
		IDToken:      "testIDToken",
	}
	tokensFetcherMock := mocks.NewThirdPartyTokensFetcher(t)
	tokensFetcherMock.On("Fetch", mock.Anything, mock.Anything).Return(testTokens, nil)
	testCase.TokensFetcherFactory.On("Create", mock.Anything).Return(tokensFetcherMock, nil)
	testEmail := "testEmail"
	testCase.TokensToEmailTransformer.On("Transform", mock.Anything).Return(testEmail, nil)
	testUUID, _ := uuid.NewUUID()
	testUser := &user.User{
		ID: testUUID,
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	responseSession := response.Content.(*session.Session)
	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if responseSession == nil {
		t.Fatal("Expected use case to return a session")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.AuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.Code, request.CallbackURL)
	testCase.TokensToEmailTransformer.AssertCalled(t, "Transform", testTokens)
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, testEmail)
}
