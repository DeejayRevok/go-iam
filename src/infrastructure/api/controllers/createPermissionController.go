package controllers

import (
	"go-uaa/src/application/createPermission"
	"go-uaa/src/domain/internals"
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
}

func (controller *CreatePermissionController) Handle(c echo.Context) error {
	accessToken, err := controller.accessTokenFinder.Find(c.Request())
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
	useCaseResponse := controller.useCaseExecutor.Execute(controller.createPermissionUseCase, &createPermissionRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusCreated)
}

func NewCreatePermissionController(useCase *createPermission.CreatePermissionUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CreatePermissionController {
	return &CreatePermissionController{
		createPermissionUseCase: useCase,
		useCaseExecutor:         useCaseExecutor,
		accessTokenFinder:       accessTokenFinder,
		dtoDeserializer:         dtoDeserializer,
		errorTransformer:        errorTransformer,
	}
}
