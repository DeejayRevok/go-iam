package controllers

import (
	"go-uaa/src/application/getAuthenticatedUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"

	"github.com/labstack/echo/v4"
)

type GetAuthenticatedUserController struct {
	getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase
	useCaseExecutor             *internals.AuthorizedUseCaseExecutor
	accessTokenFinder           *api.HTTPAccessTokenFinder
	userToResponseTransformer   *transformers.UserToResponseTransformer
	dtoSerializer               *dto.EchoDTOSerializer
	errorTransformer            *transformers.ErrorToEchoErrorTransformer
}

func (controller *GetAuthenticatedUserController) Handle(c echo.Context) error {
	accessToken, err := controller.accessTokenFinder.Find(c.Request())
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	useCaseRequest := getAuthenticatedUser.GetAuthenticatedUserRequest{
		Token: *accessToken,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(controller.getAuthenticatedUserUseCase, &useCaseRequest, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}

	userResponse := controller.userToResponseTransformer.Transform(useCaseResponse.Content.(*user.User))
	return controller.dtoSerializer.Serialize(c, userResponse)
}

func NewGetAuthenticatedUserController(useCase *getAuthenticatedUser.GetAuthenticatedUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, userTransformer *transformers.UserToResponseTransformer, dtoSerializer *dto.EchoDTOSerializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetAuthenticatedUserController {
	controller := GetAuthenticatedUserController{
		getAuthenticatedUserUseCase: useCase,
		useCaseExecutor:             useCaseExecutor,
		accessTokenFinder:           accessTokenFinder,
		userToResponseTransformer:   userTransformer,
		dtoSerializer:               dtoSerializer,
		errorTransformer:            errorTransformer,
	}
	return &controller
}
