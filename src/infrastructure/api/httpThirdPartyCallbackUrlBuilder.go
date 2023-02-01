package api

import (
	"fmt"
	"net/http"
)

const callbackURLPattern = "%s://%s%s/auth/%s/callback"

type HTTPThirdPartyCallbackURLBuilder struct {
	serverBaseProtocol string
	serverBasePath     string
}

func (builder *HTTPThirdPartyCallbackURLBuilder) Build(provider string, request *http.Request) string {
	host := request.Host
	scheme := builder.serverBaseProtocol
	return fmt.Sprintf(callbackURLPattern, scheme, host, builder.serverBasePath, provider)
}

func NewHTTPThirdPartyCallbackURLBuilder(serverBaseProtocol string, serverBasePath string) *HTTPThirdPartyCallbackURLBuilder {
	return &HTTPThirdPartyCallbackURLBuilder{
		serverBaseProtocol: serverBaseProtocol,
		serverBasePath:     serverBasePath,
	}
}
