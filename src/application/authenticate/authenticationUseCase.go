package authenticate

import (
	"context"
	"fmt"
	"go-uaa/src/domain/auth"
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/auth/authenticationStrategy"
	"go-uaa/src/domain/auth/refreshToken"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"time"
)

type AuthenticationUseCase struct {
	authenticator         *auth.Authenticator
	accessTokenGenerator  *accessToken.AccessTokenGenerator
	refreshTokenGenerator *refreshToken.RefreshTokenGenerator
	userRepository        user.UserRepository
	logger                internals.Logger
}

func (useCase *AuthenticationUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*AuthenticationRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting authentication for %s", validatedRequest.Username))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished authentication for %s", validatedRequest.Username))

	user, err := useCase.authenticator.Authenticate(ctx, validatedRequest.GrantType, useCase.createAuthenticationStrategyRequest(validatedRequest))
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	accessTokenRequest := accessToken.AccessTokenRequest{
		User:   user,
		Issuer: validatedRequest.Issuer,
	}
	accessToken := useCase.accessTokenGenerator.Generate(&accessTokenRequest)

	refreshTokenRequest := refreshToken.RefreshTokenRequest{
		User: user,
	}
	refreshToken := useCase.refreshTokenGenerator.Generate(&refreshTokenRequest)
	if err = useCase.updateUserRefreshToken(ctx, user, &refreshToken); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	authentication := useCase.createAuthentication(&accessToken, &refreshToken)
	responseContent := AuthenticationResponse{
		Authentication: &authentication,
	}
	return internals.UseCaseResponse{
		Content: &responseContent,
		Err:     nil,
	}
}

func (useCase *AuthenticationUseCase) createAuthenticationStrategyRequest(request *AuthenticationRequest) *authenticationStrategy.AuthenticationStrategyRequest {
	strategyRequest := authenticationStrategy.AuthenticationStrategyRequest{
		Username:     request.Username,
		Password:     request.Password,
		RefreshToken: request.RefreshToken,
		ThirdPartyAuthRequest: &thirdParty.ThirdPartyAuthRequest{
			State:        request.ThirdPartyState,
			Code:         request.ThirdPartyCode,
			CallbackURL:  request.ThirdPartyCallbackURL,
			AuthProvider: request.ThirdPartyAuthProvider,
		},
	}
	return &strategyRequest
}

func (useCase *AuthenticationUseCase) createAuthentication(accessTokenInstance *accessToken.AccessToken, refreshToken *refreshToken.RefreshToken) auth.Authentication {
	return auth.Authentication{
		AccessToken:  accessTokenInstance,
		RefreshToken: refreshToken,
		ExpiresIn:    int((time.Hour * time.Duration(accessToken.AccessTokenDefaultExpirationHours)).Seconds()),
		TokenType:    auth.DefaultTokenType,
	}
}

func (useCase *AuthenticationUseCase) updateUserRefreshToken(ctx context.Context, user *user.User, refreshToken *refreshToken.RefreshToken) error {
	user.RefreshToken = refreshToken.Id
	return useCase.userRepository.Save(ctx, *user)
}

func (*AuthenticationUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewAuthenticationUseCase(authenticator *auth.Authenticator, accesTokenGenerator *accessToken.AccessTokenGenerator, refreshTokenGenerator *refreshToken.RefreshTokenGenerator, userRepository user.UserRepository, logger internals.Logger) *AuthenticationUseCase {
	useCase := AuthenticationUseCase{
		authenticator:         authenticator,
		accessTokenGenerator:  accesTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
		userRepository:        userRepository,
		logger:                logger,
	}
	return &useCase
}
