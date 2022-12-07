package authenticate

import (
	"context"
	"go-iam/src/domain/auth"
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/auth/authenticationStrategy"
	"go-iam/src/domain/auth/refreshToken"
	"go-iam/src/domain/auth/thirdParty"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
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

	useCase.logger.Info(ctx, "Starting authentication")
	defer useCase.logger.Info(ctx, "Finished authentication")

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
		Email:        request.Email,
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
