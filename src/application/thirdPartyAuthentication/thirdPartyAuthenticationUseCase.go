package thirdPartyAuthentication

import (
	"context"
	"fmt"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/session"
	"go-uaa/src/domain/user"
)

type ThirdPartyAuthenticationUseCase struct {
	authStateChecker         *thirdParty.ThirdPartyAuthStateChecker
	tokensFetcherFactory     thirdParty.ThirdPartyTokensFetcherFactory
	tokensToEmailTransformer thirdParty.ThirdPartyTokensToEmailTransformer
	userRepository           user.UserRepository
	sessionGenerator         *session.SessionGenerator
	logger                   internals.Logger
}

func (useCase *ThirdPartyAuthenticationUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*ThirdPartyAuthenticationRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting authentication for %s", validatedRequest.AuthProvider))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished authentication for %s", validatedRequest.AuthProvider))

	if err := useCase.authStateChecker.Check(validatedRequest.State); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	tokensFetcher, err := useCase.tokensFetcherFactory.Create(validatedRequest.AuthProvider)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	tokens, err := tokensFetcher.Fetch(validatedRequest.Code, validatedRequest.CallbackURL)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	email, err := useCase.tokensToEmailTransformer.Transform(tokens)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	user, err := useCase.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if user == nil {
		return internals.ErrorUseCaseResponse(fmt.Errorf("user with email %s not registered", email))
	}
	session := useCase.sessionGenerator.Generate(user.ID)
	return internals.UseCaseResponse{
		Content: session,
		Err:     nil,
	}
}

func (*ThirdPartyAuthenticationUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewThirdPartyAuthenticationUseCase(authStateChecker *thirdParty.ThirdPartyAuthStateChecker, tokensFetcherFactory thirdParty.ThirdPartyTokensFetcherFactory, tokensToEmailTransformer thirdParty.ThirdPartyTokensToEmailTransformer, userRepository user.UserRepository, logger internals.Logger) *ThirdPartyAuthenticationUseCase {
	return &ThirdPartyAuthenticationUseCase{
		authStateChecker:         authStateChecker,
		tokensFetcherFactory:     tokensFetcherFactory,
		tokensToEmailTransformer: tokensToEmailTransformer,
		userRepository:           userRepository,
		logger:                   logger,
	}
}
