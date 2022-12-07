package controllers

import (
	"go-iam/src/application/getThirdPartyAuthenticationUrl"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/api"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetThirdPartyAuthenticationController struct {
	callbackURLBuilder                    *api.HTTPThirdPartyCallbackURLBuilder
	getThirdPartyAuthenticationUrlUseCase *getThirdPartyAuthenticationUrl.GetThirdPartyAuthenticationURLUseCase
	useCaseExecutor                       *internals.UseCaseExecutor
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

	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.getThirdPartyAuthenticationUrlUseCase, &getUrlRequest, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.Redirect(http.StatusTemporaryRedirect, useCaseResponse.Content.(string))
}

func NewGetThirdPartyAuthenticationController(callbackURLBuilder *api.HTTPThirdPartyCallbackURLBuilder, getThirdPartyAuthenticationUrlUseCase *getThirdPartyAuthenticationUrl.GetThirdPartyAuthenticationURLUseCase, useCaseExecutor *internals.UseCaseExecutor, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetThirdPartyAuthenticationController {
	return &GetThirdPartyAuthenticationController{
		callbackURLBuilder:                    callbackURLBuilder,
		getThirdPartyAuthenticationUrlUseCase: getThirdPartyAuthenticationUrlUseCase,
		useCaseExecutor:                       useCaseExecutor,
		errorTransformer:                      errorTransformer,
	}
}
