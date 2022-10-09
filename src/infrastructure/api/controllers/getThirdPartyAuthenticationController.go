package controllers

import (
	"go-uaa/src/application/getThirdPartyAuthenticationUrl"
	"go-uaa/src/domain/internals"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetThirdPartyAuthenticationController struct {
	callbackURLBuilder                    *api.HTTPThirdPartyCallbackURLBuilder
	getThirdPartyAuthenticationUrlUseCase *getThirdPartyAuthenticationUrl.GetThirdPartyAuthenticationURLUseCase
	useCaseExecutor                       *internals.AuthorizedUseCaseExecutor
	errorTransformer                      *transformers.ErrorToEchoErrorTransformer
}

func (controller *GetThirdPartyAuthenticationController) Handle(c echo.Context) error {
	authProvider := c.Param("provider")
	if authProvider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authentication provider")
	}

	request := c.Request()
	ctx := request.Context()
	getUrlRequest := getThirdPartyAuthenticationUrl.GetThirdPartyAuthenticationURLRequest{
		AuthProvider: authProvider,
		CallbackURL:  controller.callbackURLBuilder.Build(authProvider, request),
	}

	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.getThirdPartyAuthenticationUrlUseCase, &getUrlRequest, nil, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.Redirect(http.StatusTemporaryRedirect, useCaseResponse.Content.(string))
}

func NewGetThirdPartyAuthenticationController(callbackURLBuilder *api.HTTPThirdPartyCallbackURLBuilder, getThirdPartyAuthenticationUrlUseCase *getThirdPartyAuthenticationUrl.GetThirdPartyAuthenticationURLUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetThirdPartyAuthenticationController {
	return &GetThirdPartyAuthenticationController{
		callbackURLBuilder:                    callbackURLBuilder,
		getThirdPartyAuthenticationUrlUseCase: getThirdPartyAuthenticationUrlUseCase,
		useCaseExecutor:                       useCaseExecutor,
		errorTransformer:                      errorTransformer,
	}
}
