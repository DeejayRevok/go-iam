package api

import (
	"go-iam/src/domain/auth/accessToken"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type HTTPAccessTokenFinder struct {
	tokenDeserializer accessToken.AccessTokenDeserializer
}

func (finder *HTTPAccessTokenFinder) Find(httpRequest *http.Request) (*accessToken.AccessToken, error) {
	serializedToken, err := finder.getSerializedAccessToken(*httpRequest)
	if err != nil {
		return nil, err
	}
	if serializedToken == "" {
		return nil, nil
	}
	token, err := finder.tokenDeserializer.Deserialize(serializedToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (*HTTPAccessTokenFinder) getSerializedAccessToken(request http.Request) (string, error) {
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", nil
	}
	splittedHeader := strings.Split(authorizationHeader, " ")
	if len(splittedHeader) < 2 {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Malformed authorization header")
	}
	return splittedHeader[1], nil
}

func NewHTTPAccessTokenFinder(tokenDeserializer accessToken.AccessTokenDeserializer) *HTTPAccessTokenFinder {
	return &HTTPAccessTokenFinder{
		tokenDeserializer: tokenDeserializer,
	}
}
