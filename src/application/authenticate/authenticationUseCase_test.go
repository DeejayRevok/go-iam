package authenticate

import (
	"context"
	"fmt"
	"go-uaa/mocks"
	"go-uaa/src/domain/auth"
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/auth/authenticationStrategy"
	"go-uaa/src/domain/auth/refreshToken"
	"go-uaa/src/domain/session"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepository           *mocks.UserRepository
	HashComparator           *mocks.HashComparator
	RefreshTokenDeserializer *mocks.RefreshTokenDeserializer
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
	authenticator := auth.NewAuthenticator(passwordAuthStrategy, refreshTokenAuthStrategy)
	accessTokenGenerator := accessToken.NewAccessTokenGenerator()
	refreshTokenGenerator := refreshToken.NewRefreshTokenGenerator()
	sessionGenerator := session.NewSessionGenerator()
	return testCase{
		UserRepository:           userRepository,
		HashComparator:           hashComparator,
		RefreshTokenDeserializer: refreshTokenDeserializer,
		UseCase:                  NewAuthenticationUseCase(authenticator, accessTokenGenerator, refreshTokenGenerator, userRepository, sessionGenerator, logger),
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
		Username:     "test",
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
		Username:     "testUsername",
		Password:     "testPassword",
		Issuer:       "testIssuer",
		GrantType:    "password",
		RefreshToken: "",
	}
	testUser := user.User{
		Username:     "testUsername",
		Password:     "testPassword",
		RefreshToken: "",
	}
	testCase.UserRepository.On("FindByUsername", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	authenticateResponse := *response.Content.(*AuthenticationResponse)
	responseAuthentication := authenticateResponse.Authentication
	responseSession := authenticateResponse.Session

	if responseAuthentication.AccessToken == nil {
		t.Fatal("Expected use case to return an access token")
	}
	if responseAuthentication.AccessToken.Iss != request.Issuer {
		t.Fatal("Expected use case to return an access token with same issuer than the request")
	}
	if responseAuthentication.AccessToken.Sub != testUser.Username {
		t.Fatal("Expected use case to return an access token with sub equals to the found user username")
	}
	if responseAuthentication.RefreshToken.Sub != testUser.Username {
		t.Fatal("Expected use case to return a refresh token with sub equals to the found user username")
	}
	if responseAuthentication.ExpiresIn != 3600 {
		t.Fatal("Expected use case to return one hour of expiration")
	}
	if responseAuthentication.TokenType != "bearer" {
		t.Fatal("Expected use case to return bearer token type")
	}
	if responseSession == nil {
		t.Fatal("Expected use to return a session")
	}
	if responseSession.UserID != testUser.ID.String() {
		t.Fatal("Expected session user id to be the same as the found user id")
	}
	testCase.UserRepository.AssertCalled(t, "FindByUsername", ctx, request.Username)
	testCase.HashComparator.AssertCalled(t, "Compare", request.Password, testUser.Password)
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUser.Username && user.RefreshToken != ""
	}))
	testCase.RefreshTokenDeserializer.AssertNotCalled(t, "Deserialize")
}

func TestExecuteRefreshTokenGrantTypeSuccess(t *testing.T) {
	testCase := setUp(t)
	request := AuthenticationRequest{
		Username:     "testUsername",
		Password:     "testPassword",
		Issuer:       "testIssuer",
		GrantType:    "refresh_token",
		RefreshToken: "testRefreshTokenId",
	}
	ctx := context.Background()
	testRefreshToken := refreshToken.RefreshToken{
		Id:  "testRefreshTokenId",
		Sub: "testUsername",
		Exp: 3600,
	}
	testUser := user.User{
		Username:     "testUsername",
		Password:     "testPassword",
		RefreshToken: "testRefreshTokenId",
	}
	testCase.UserRepository.On("FindByUsername", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.RefreshTokenDeserializer.On("Deserialize", mock.Anything).Return(&testRefreshToken, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatalf("Expected use case not to return error. Returned: %s", response.Err.Error())
	}
	authenticateResponse := *response.Content.(*AuthenticationResponse)
	responseAuthentication := authenticateResponse.Authentication
	responseSession := authenticateResponse.Session
	if responseAuthentication.AccessToken == nil {
		t.Fatal("Expected use case to return an access token")
	}
	if responseAuthentication.AccessToken.Iss != request.Issuer {
		t.Fatal("Expected use case to return an access token with same issuer than the request")
	}
	if responseAuthentication.AccessToken.Sub != testUser.Username {
		t.Fatal("Expected use case to return an access token with sub equals to the found user username")
	}
	if responseAuthentication.RefreshToken.Sub != testUser.Username {
		t.Fatal("Expected use case to return a refresh token with sub equals to the found user username")
	}
	if responseAuthentication.ExpiresIn != 3600 {
		t.Fatal("Expected use case to return one hour of expiration")
	}
	if responseAuthentication.TokenType != "bearer" {
		t.Fatal("Expected use case to return bearer token type")
	}
	if responseSession == nil {
		t.Fatal("Expected use to return a session")
	}
	if responseSession.UserID != testUser.ID.String() {
		t.Fatal("Expected session user id to be the same as the found user id")
	}
	testCase.UserRepository.AssertCalled(t, "FindByUsername", ctx, request.Username)
	testCase.HashComparator.AssertNotCalled(t, "Compare")
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUser.Username && user.RefreshToken != ""
	}))
	testCase.RefreshTokenDeserializer.AssertCalled(t, "Deserialize", request.RefreshToken)
}
