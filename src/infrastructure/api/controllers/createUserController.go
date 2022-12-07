package controllers

import (
	"go-iam/src/application/createUser"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserController struct {
	createUserUseCase *createUser.CreateUserUseCase
	useCaseExecutor   *internals.UseCaseExecutor
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
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.createUserUseCase, &createUserRequest, nil)
	if err := controller.handleUseCaseError(useCaseResponse.Err); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
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

func NewCreateUserController(createUserUseCase *createUser.CreateUserUseCase, useCaseExecutor *internals.UseCaseExecutor, echoDTODeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CreateUserController {
	controller := CreateUserController{
		createUserUseCase: createUserUseCase,
		useCaseExecutor:   useCaseExecutor,
		dtoDeserializer:   echoDTODeserializer,
		errorTransformer:  errorTransformer,
	}
	return &controller
}
