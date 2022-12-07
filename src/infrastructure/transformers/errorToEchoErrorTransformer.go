package transformers

import (
	"go-iam/src/domain/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorToEchoErrorTransformer struct{}

func (transformer *ErrorToEchoErrorTransformer) Transform(err error) *echo.HTTPError {
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(transformer.getHTTPStatusCode(err), err.Error())
}

func (*ErrorToEchoErrorTransformer) getHTTPStatusCode(err error) int {
	switch err.(type) {
	case auth.AuthenticationError:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func NewErrorToEchoErrorTransformer() *ErrorToEchoErrorTransformer {
	return &ErrorToEchoErrorTransformer{}
}
