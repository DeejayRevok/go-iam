package controllers

import (
	"go-uaa/src/application/createPermission"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/session"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreatePermissionController struct {
	createPermissionUseCase *createPermission.CreatePermissionUseCase
	useCaseExecutor         *internals.AuthorizedUseCaseExecutor
	accessTokenFinder       *api.HTTPAccessTokenFinder
	dtoDeserializer         *dto.EchoDTODeserializer
	errorTransformer        *transformers.ErrorToEchoErrorTransformer
	sessionFinder           *api.HTTPSessionFinder
}

func (controller *CreatePermissionController) Handle(c echo.Context) error {
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	var requestSession *session.Session
	if accessToken == nil {
		requestSession, err = controller.sessionFinder.Find(request)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var creationRequestDTO dto.PermissionCreationRequestDTO
	if err := controller.dtoDeserializer.Deserialize(c, &creationRequestDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	createPermissionRequest := createPermission.CreatePermissionRequest{
		Name: creationRequestDTO.Name,
	}
	ctx := c.Request().Context()
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.createPermissionUseCase, &createPermissionRequest, accessToken, requestSession)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusCreated)
}

func NewCreatePermissionController(useCase *createPermission.CreatePermissionUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer, sessionFinder *api.HTTPSessionFinder) *CreatePermissionController {
	return &CreatePermissionController{
		createPermissionUseCase: useCase,
		useCaseExecutor:         useCaseExecutor,
		accessTokenFinder:       accessTokenFinder,
		dtoDeserializer:         dtoDeserializer,
		errorTransformer:        errorTransformer,
		sessionFinder:           sessionFinder,
	}
}
