package controllers

import (
	"go-iam/src/application/authenticate"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthenticateController struct {
	authenticateUseCase       *authenticate.AuthenticationUseCase
	useCaseExecutor           *internals.UseCaseExecutor
	authenticationTransformer *transformers.AuthenticationToResponseTransformer
	dtoDeserializer           *dto.EchoDTODeserializer
	dtoSerializer             *dto.EchoDTOSerializer
	errorTransformer          *transformers.ErrorToEchoErrorTransformer
}

func (controller *AuthenticateController) Handle(c echo.Context) error {
	var authRequestDTO dto.AuthenticationRequestDTO
	err := controller.dtoDeserializer.Deserialize(c, &authRequestDTO)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	httpRequest := c.Request()
	authenticationRequest := authenticate.AuthenticationRequest{
		Email:        authRequestDTO.Email,
		Password:     authRequestDTO.Password,
		Issuer:       controller.getRequestOrigin(httpRequest),
		GrantType:    authRequestDTO.GrantType,
		RefreshToken: authRequestDTO.RefreshToken,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(httpRequest.Context(), controller.authenticateUseCase, &authenticationRequest, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	authenticationResponse := useCaseResponse.Content.(*authenticate.AuthenticationResponse)
	authenticationDTO, err := controller.authenticationTransformer.Transform(authenticationResponse.Authentication)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	return controller.dtoSerializer.Serialize(c, authenticationDTO)
}

func (controller *AuthenticateController) getRequestOrigin(request *http.Request) string {
	return request.Header.Get("origin")
}

func NewAuthenticateController(useCase *authenticate.AuthenticationUseCase, useCaseExecutor *internals.UseCaseExecutor, transformer *transformers.AuthenticationToResponseTransformer, dtoDeserializer *dto.EchoDTODeserializer, dtoSerializer *dto.EchoDTOSerializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *AuthenticateController {
	return &AuthenticateController{
		authenticateUseCase:       useCase,
		useCaseExecutor:           useCaseExecutor,
		authenticationTransformer: transformer,
		dtoDeserializer:           dtoDeserializer,
		dtoSerializer:             dtoSerializer,
		errorTransformer:          errorTransformer,
	}
}
