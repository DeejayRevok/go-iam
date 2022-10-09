package controllers

import (
	"go-uaa/src/application/createUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserController struct {
	createUserUseCase *createUser.CreateUserUseCase
	useCaseExecutor   *internals.AuthorizedUseCaseExecutor
	dtoDeserializer   *dto.EchoDTODeserializer
	errorTransformer  *transformers.ErrorToEchoErrorTransformer
}

func (controller *CreateUserController) Handle(c echo.Context) error {
	var userCreationRequest dto.UserCreationRequestDTO
	if err := controller.dtoDeserializer.Deserialize(c, &userCreationRequest); err != nil {
		return controller.errorTransformer.Transform(err)
	}

	ctx := c.Request().Context()
	createUserRequest := createUser.CreateUserRequest{
		Username: *userCreationRequest.Username,
		Email:    *userCreationRequest.Email,
		Password: *userCreationRequest.Password,
		Roles:    controller.parseRoles(*userCreationRequest.Roles),
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.createUserUseCase, &createUserRequest, nil, nil)
	if err := controller.handleUseCaseError(useCaseResponse.Err); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (resolver *CreateUserController) parseRoles(roles []*string) []string {
	parsedRoles := make([]string, 0)
	for _, role := range roles {
		parsedRoles = append(parsedRoles, *role)
	}
	return parsedRoles
}

func (controller *CreateUserController) handleUseCaseError(err error) error {
	if err != nil {
		switch err.(type) {
		case user.EmailNotValidError:
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}

func NewCreateUserController(createUserUseCase *createUser.CreateUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, echoDTODeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CreateUserController {
	controller := CreateUserController{
		createUserUseCase: createUserUseCase,
		useCaseExecutor:   useCaseExecutor,
		dtoDeserializer:   echoDTODeserializer,
		errorTransformer:  errorTransformer,
	}
	return &controller
}
