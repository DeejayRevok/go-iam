package controllers

import (
	"go-iam/src/application/getAuthenticatedUser"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/api"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/transformers"

	"github.com/labstack/echo/v4"
)

type GetAuthenticatedUserController struct {
	getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase
	useCaseExecutor             *internals.UseCaseExecutor
	accessTokenFinder           *api.HTTPAccessTokenFinder
	userToResponseTransformer   *transformers.UserToResponseTransformer
	dtoSerializer               *dto.EchoDTOSerializer
	errorTransformer            *transformers.ErrorToEchoErrorTransformer
}

func (controller *GetAuthenticatedUserController) Handle(c echo.Context) error {
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	ctx := request.Context()
	useCaseRequest := getAuthenticatedUser.GetAuthenticatedUserRequest{
		Token: accessToken,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.getAuthenticatedUserUseCase, &useCaseRequest, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}

	userResponse := controller.userToResponseTransformer.Transform(useCaseResponse.Content.(*user.User))
	return controller.dtoSerializer.Serialize(c, userResponse)
}

func NewGetAuthenticatedUserController(useCase *getAuthenticatedUser.GetAuthenticatedUserUseCase, useCaseExecutor *internals.UseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, userTransformer *transformers.UserToResponseTransformer, dtoSerializer *dto.EchoDTOSerializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetAuthenticatedUserController {
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
