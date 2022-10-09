package getThirdPartyAuthenticationUrl

import (
	"context"
	"fmt"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/domain/internals"
)

type GetThirdPartyAuthenticationURLUseCase struct {
	authURLBuilderFactory thirdParty.ThirdPartyAuthURLBuilderFactory
	logger                internals.Logger
}

func (useCase *GetThirdPartyAuthenticationURLUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*GetThirdPartyAuthenticationURLRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting getting auth url for %s", validatedRequest.AuthProvider))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished getting auth url for %s", validatedRequest.AuthProvider))

	authURLBuilder, err := useCase.authURLBuilderFactory.Create(validatedRequest.AuthProvider)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	authURL := authURLBuilder.Build(validatedRequest.CallbackURL)
	return internals.UseCaseResponse{
		Content: authURL,
		Err:     nil,
	}
}

func (*GetThirdPartyAuthenticationURLUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewGetThirdPartyAuthenticationURLUseCase(authURLBuilderFactory thirdParty.ThirdPartyAuthURLBuilderFactory, logger internals.Logger) *GetThirdPartyAuthenticationURLUseCase {
	return &GetThirdPartyAuthenticationURLUseCase{
		authURLBuilderFactory: authURLBuilderFactory,
		logger:                logger,
	}
}
