package authenticate

import (
	"context"
	"fmt"
	"go-iam/mocks"
	"go-iam/src/domain/auth"
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/auth/authenticationStrategy"
	"go-iam/src/domain/auth/refreshToken"
	"go-iam/src/domain/auth/thirdParty"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepository           *mocks.UserRepository
	HashComparator           *mocks.HashComparator
	RefreshTokenDeserializer *mocks.RefreshTokenDeserializer
	TokensFetcherFactory     *mocks.ThirdPartyTokensFetcherFactory
	TokensToEmailTransformer *mocks.ThirdPartyTokensToEmailTransformer
	InternalState            string
	UseCase                  *AuthenticationUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepository := mocks.NewUserRepository(t)
	hashComparator := mocks.NewHashComparator(t)
	refreshTokenDeserializer := mocks.NewRefreshTokenDeserializer(t)
	passwordAuthStrategy := authenticationStrategy.NewPasswordAuthenticationStrategy(userRepository, hashComparator)
	refreshTokenAuthStrategy := authenticationStrategy.NewRefreshTokenAuthenticationStrategy(userRepository, refreshTokenDeserializer)
	tokensFetcherFactoryMock := mocks.NewThirdPartyTokensFetcherFactory(t)
	tokensToEmailTransformerMock := mocks.NewThirdPartyTokensToEmailTransformer(t)
	testInternalState := "testInternalState"
	authStateChecker := thirdParty.NewThirdPartyAuthStateChecker(testInternalState)
	thirdPartyAuthStrategy := authenticationStrategy.NewThirdPartyAuthenticationStrategy(authStateChecker, tokensFetcherFactoryMock, tokensToEmailTransformerMock, userRepository)
	authenticator := auth.NewAuthenticator(passwordAuthStrategy, refreshTokenAuthStrategy, thirdPartyAuthStrategy)
	accessTokenGenerator := accessToken.NewAccessTokenGenerator()
	refreshTokenGenerator := refreshToken.NewRefreshTokenGenerator()
	return testCase{
		UserRepository:           userRepository,
		HashComparator:           hashComparator,
		RefreshTokenDeserializer: refreshTokenDeserializer,
		TokensFetcherFactory:     tokensFetcherFactoryMock,
		TokensToEmailTransformer: tokensToEmailTransformerMock,
		InternalState:            testInternalState,
		UseCase:                  NewAuthenticationUseCase(authenticator, accessTokenGenerator, refreshTokenGenerator, userRepository, logger),
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
	testCase.UserRepository.AssertNotCalled(t, "Save")
	testCase.HashComparator.AssertNotCalled(t, "Compare")
	testCase.RefreshTokenDeserializer.AssertNotCalled(t, "Deserialize")
}

func TestExecuteAuthenticatorFailsWrongGrantType(t *testing.T) {
	testCase := setUp(t)
	wrongGrantType := "wrong"
	request := AuthenticationRequest{
		Email:        "test",
		Password:     "test",
		Issuer:       "test",
		GrantType:    wrongGrantType,
		RefreshToken: "test",
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err.Error() != fmt.Sprintf("grant type %s is not supported", wrongGrantType) {
		t.Fatal("Expected use case to return grant type not supported error")
	}
	testCase.UserRepository.AssertNotCalled(t, "Save")
	testCase.HashComparator.AssertNotCalled(t, "Compare")
	testCase.RefreshTokenDeserializer.AssertNotCalled(t, "Deserialize")
}

func TestExecutePasswordGrantTypeSuccess(t *testing.T) {
	testCase := setUp(t)
	request := AuthenticationRequest{
		Email:        "testUsername",
		Password:     "testPassword",
		Issuer:       "testIssuer",
		GrantType:    "password",
		RefreshToken: "",
	}
	testUser := user.User{
		Username:     "testUsername",
		Password:     "testPassword",
		Email:        "testEmail",
		RefreshToken: "",
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	authenticateResponse := *response.Content.(*AuthenticationResponse)
	responseAuthentication := authenticateResponse.Authentication

	if responseAuthentication.AccessToken == nil {
		t.Fatal("Expected use case to return an access token")
	}
	if responseAuthentication.AccessToken.Iss != request.Issuer {
		t.Fatal("Expected use case to return an access token with same issuer than the request")
	}
	if responseAuthentication.AccessToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return an access token with sub equals to the found user email")
	}
	if responseAuthentication.RefreshToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return a refresh token with sub equals to the found user email")
	}
	if responseAuthentication.ExpiresIn != 3600 {
		t.Fatal("Expected use case to return one hour of expiration")
	}
	if responseAuthentication.TokenType != "bearer" {
		t.Fatal("Expected use case to return bearer token type")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, request.Email)
	testCase.HashComparator.AssertCalled(t, "Compare", request.Password, testUser.Password)
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUser.Username && user.RefreshToken != ""
	}))
	testCase.RefreshTokenDeserializer.AssertNotCalled(t, "Deserialize")
}

func TestExecuteRefreshTokenGrantTypeSuccess(t *testing.T) {
	testCase := setUp(t)
	request := AuthenticationRequest{
		Email:        "testEmail",
		Password:     "testPassword",
		Issuer:       "testIssuer",
		GrantType:    "refresh_token",
		RefreshToken: "testRefreshTokenId",
	}
	ctx := context.Background()
	testRefreshToken := refreshToken.RefreshToken{
		Id:  "testRefreshTokenId",
		Sub: "testEmail",
		Exp: 3600,
	}
	testUser := user.User{
		Username:     "testUsername",
		Email:        "testEmail",
		Password:     "testPassword",
		RefreshToken: "testRefreshTokenId",
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.RefreshTokenDeserializer.On("Deserialize", mock.Anything).Return(&testRefreshToken, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatalf("Expected use case not to return error. Returned: %s", response.Err.Error())
	}
	authenticateResponse := *response.Content.(*AuthenticationResponse)
	responseAuthentication := authenticateResponse.Authentication
	if responseAuthentication.AccessToken == nil {
		t.Fatal("Expected use case to return an access token")
	}
	if responseAuthentication.AccessToken.Iss != request.Issuer {
		t.Fatal("Expected use case to return an access token with same issuer than the request")
	}
	if responseAuthentication.AccessToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return an access token with sub equals to the found user email")
	}
	if responseAuthentication.RefreshToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return a refresh token with sub equals to the found user email")
	}
	if responseAuthentication.ExpiresIn != 3600 {
		t.Fatal("Expected use case to return one hour of expiration")
	}
	if responseAuthentication.TokenType != "bearer" {
		t.Fatal("Expected use case to return bearer token type")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, request.Email)
	testCase.HashComparator.AssertNotCalled(t, "Compare")
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUser.Username && user.RefreshToken != ""
	}))
	testCase.RefreshTokenDeserializer.AssertCalled(t, "Deserialize", request.RefreshToken)
}

func TestThirdPartyGrantTypeSuccess(t *testing.T) {
	testCase := setUp(t)
	request := AuthenticationRequest{
		Issuer:                 "testIssuer",
		GrantType:              "third_party",
		ThirdPartyState:        testCase.InternalState,
		ThirdPartyCode:         "testCode",
		ThirdPartyAuthProvider: "testAuthProvider",
		ThirdPartyCallbackURL:  "testCallbackURL",
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
	testUser := user.User{
		Username:     "testUsername",
		Password:     "testPassword",
		RefreshToken: "testRefreshTokenId",
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatalf("Expected use case not to return error. Returned: %s", response.Err.Error())
	}
	authenticateResponse := *response.Content.(*AuthenticationResponse)
	responseAuthentication := authenticateResponse.Authentication
	if responseAuthentication.AccessToken == nil {
		t.Fatal("Expected use case to return an access token")
	}
	if responseAuthentication.AccessToken.Iss != request.Issuer {
		t.Fatal("Expected use case to return an access token with same issuer than the request")
	}
	if responseAuthentication.AccessToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return an access token with sub equals to the found user email")
	}
	if responseAuthentication.RefreshToken.Sub != testUser.Email {
		t.Fatal("Expected use case to return a refresh token with sub equals to the found user email")
	}
	if responseAuthentication.ExpiresIn != 3600 {
		t.Fatal("Expected use case to return one hour of expiration")
	}
	if responseAuthentication.TokenType != "bearer" {
		t.Fatal("Expected use case to return bearer token type")
	}
	testCase.TokensFetcherFactory.AssertCalled(t, "Create", request.ThirdPartyAuthProvider)
	tokensFetcherMock.AssertCalled(t, "Fetch", request.ThirdPartyCode, request.ThirdPartyCallbackURL)
	testCase.TokensToEmailTransformer.AssertCalled(t, "Transform", testTokens)
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUser.Username && user.RefreshToken != ""
	}))
}
