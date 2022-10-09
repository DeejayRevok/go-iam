package controllers

import (
	"go-uaa/src/application/thirdPartyAuthentication"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/session"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ThirdPartyAuthenticationCallbackController struct {
	callbackURLBuilder              *api.HTTPThirdPartyCallbackURLBuilder
	thirdPartyAuthenticationUseCase *thirdPartyAuthentication.ThirdPartyAuthenticationUseCase
	useCaseExecutor                 *internals.AuthorizedUseCaseExecutor
	errorTransformer                *transformers.ErrorToEchoErrorTransformer
	sessionSetter                   *api.EchoSessionSetter
}

func (controller *ThirdPartyAuthenticationCallbackController) Handle(c echo.Context) error {
	authProvider := c.Param("provider")
	if authProvider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authentication provider")
	}

	state := c.FormValue("state")
	if state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authentication state")
	}

	code := c.FormValue("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authentication code")
	}

	request := c.Request()
	ctx := request.Context()
	authRequest := thirdPartyAuthentication.ThirdPartyAuthenticationRequest{
		State:        state,
		Code:         code,
		AuthProvider: authProvider,
		CallbackURL:  controller.callbackURLBuilder.Build(authProvider, request),
	}

	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.thirdPartyAuthenticationUseCase, &authRequest, nil, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}

	c, err := controller.sessionSetter.Set(c, useCaseResponse.Content.(*session.Session))
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	return c.NoContent(http.StatusOK)
}

func NewThirdPartyAuthenticationCallbackController(callbackURLBuilder *api.HTTPThirdPartyCallbackURLBuilder, thirdPartyAuthenticationUseCase *thirdPartyAuthentication.ThirdPartyAuthenticationUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, errorTransformer *transformers.ErrorToEchoErrorTransformer, sessionSetter *api.EchoSessionSetter) *ThirdPartyAuthenticationCallbackController {
	return &ThirdPartyAuthenticationCallbackController{
		callbackURLBuilder:              callbackURLBuilder,
		thirdPartyAuthenticationUseCase: thirdPartyAuthenticationUseCase,
		useCaseExecutor:                 useCaseExecutor,
		errorTransformer:                errorTransformer,
		sessionSetter:                   sessionSetter,
	}
}
