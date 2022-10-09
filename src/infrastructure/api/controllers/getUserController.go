package controllers

import (
	"go-uaa/src/application/getUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/session"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetUserController struct {
	getUserUseCase    *getUser.GetUserUseCase
	useCaseExecutor   *internals.AuthorizedUseCaseExecutor
	accessTokenFinder *api.HTTPAccessTokenFinder
	dtoSerializer     *dto.EchoDTOSerializer
	userTransformer   *transformers.UserToResponseTransformer
	errorTransformer  *transformers.ErrorToEchoErrorTransformer
	sessionFinder     *api.HTTPSessionFinder
}

func (controller *GetUserController) Handle(c echo.Context) error {
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

	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error processing user id")
	}
	ctx := c.Request().Context()
	userRequest := getUser.GetUserRequest{
		UserId: userId,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.getUserUseCase, &userRequest, accessToken, requestSession)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	userResponse := controller.userTransformer.Transform(useCaseResponse.Content.(*user.User))
	return controller.errorTransformer.Transform(controller.dtoSerializer.Serialize(c, userResponse))
}

func NewGetUserController(getUserUseCase *getUser.GetUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoSerializer *dto.EchoDTOSerializer, userTransformer *transformers.UserToResponseTransformer, errorTransformer *transformers.ErrorToEchoErrorTransformer, sessionFinder *api.HTTPSessionFinder) *GetUserController {
	controller := GetUserController{
		getUserUseCase:    getUserUseCase,
		useCaseExecutor:   useCaseExecutor,
		accessTokenFinder: accessTokenFinder,
		dtoSerializer:     dtoSerializer,
		userTransformer:   userTransformer,
		errorTransformer:  errorTransformer,
		sessionFinder:     sessionFinder,
	}
	return &controller
}
