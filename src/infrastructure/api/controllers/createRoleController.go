package controllers

import (
	"go-uaa/src/application/createRole"
	"go-uaa/src/domain/internals"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateRoleController struct {
	createRoleUseCase *createRole.CreateRoleUseCase
	useCaseExecutor   *internals.AuthorizedUseCaseExecutor
	accessTokenFinder *api.HTTPAccessTokenFinder
	dtoDeserializer   *dto.EchoDTODeserializer
	errorTransformer  *transformers.ErrorToEchoErrorTransformer
}

func (controller *CreateRoleController) Handle(c echo.Context) error {
	accessToken, err := controller.accessTokenFinder.Find(c.Request())
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var creationRequestDTO dto.RoleCreationRequestDTO
	if err := controller.dtoDeserializer.Deserialize(c, &creationRequestDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	ctx := c.Request().Context()
	createRoleRequest := createRole.CreateRoleRequest{
		Name:        creationRequestDTO.Name,
		Permissions: creationRequestDTO.Permissions,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.createRoleUseCase, &createRoleRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusCreated)
}

func NewCreateRoleController(useCase *createRole.CreateRoleUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CreateRoleController {
	return &CreateRoleController{
		createRoleUseCase: useCase,
		useCaseExecutor:   useCaseExecutor,
		accessTokenFinder: accessTokenFinder,
		dtoDeserializer:   dtoDeserializer,
		errorTransformer:  errorTransformer,
	}
}
