package controllers

import (
	"go-uaa/src/application/authenticate"
	"go-uaa/src/domain/internals"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ThirdPartyAuthenticationCallbackController struct {
	callbackURLBuilder        *api.HTTPThirdPartyCallbackURLBuilder
	authenticationUseCase     *authenticate.AuthenticationUseCase
	useCaseExecutor           *internals.AuthorizedUseCaseExecutor
	errorTransformer          *transformers.ErrorToEchoErrorTransformer
	authenticationTransformer *transformers.AuthenticationToResponseTransformer
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
	authRequest := authenticate.AuthenticationRequest{
		Issuer:                 controller.getRequestOrigin(request),
		ThirdPartyState:        state,
		ThirdPartyCode:         code,
		ThirdPartyAuthProvider: authProvider,
		ThirdPartyCallbackURL:  controller.callbackURLBuilder.Build(authProvider, request),
		GrantType:              "third_party",
	}

	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.authenticationUseCase, &authRequest, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	authenticationResponse := useCaseResponse.Content.(*authenticate.AuthenticationResponse)
	authenticationDTO, err := controller.authenticationTransformer.Transform(authenticationResponse.Authentication)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	controller.setAuthenticationCookies(c, authenticationDTO.AccessToken, authenticationDTO.RefreshToken)
	return c.NoContent(http.StatusOK)
}

func (controller *ThirdPartyAuthenticationCallbackController) setAuthenticationCookies(c echo.Context, accessToken string, refreshToken string) {
	controller.setCookie(c, "access_token", accessToken)
	controller.setCookie(c, "refresh_token", refreshToken)
}

func (*ThirdPartyAuthenticationCallbackController) setCookie(c echo.Context, cookieName string, cookieValue string) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = cookieValue
	cookie.HttpOnly = false
	cookie.Secure = true
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func (*ThirdPartyAuthenticationCallbackController) getRequestOrigin(request *http.Request) string {
	return request.Header.Get("origin")
}

func NewThirdPartyAuthenticationCallbackController(callbackURLBuilder *api.HTTPThirdPartyCallbackURLBuilder, authenticationUseCase *authenticate.AuthenticationUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, errorTransformer *transformers.ErrorToEchoErrorTransformer, authenticationTransformer *transformers.AuthenticationToResponseTransformer) *ThirdPartyAuthenticationCallbackController {
	return &ThirdPartyAuthenticationCallbackController{
		callbackURLBuilder:    callbackURLBuilder,
		authenticationUseCase: authenticationUseCase,
		useCaseExecutor:       useCaseExecutor,
		errorTransformer:      errorTransformer,
	}
}
