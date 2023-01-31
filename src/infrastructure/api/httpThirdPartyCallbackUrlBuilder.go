package api

import (
	"fmt"
	"net/http"
)

const callbackURLPattern = "%s://%s%s/auth/%s/callback"

type HTTPThirdPartyCallbackURLBuilder struct {
	serverBasePath string
}

func (builder *HTTPThirdPartyCallbackURLBuilder) Build(provider string, request *http.Request) string {
	host := request.Host
	scheme := builder.getRequestScheme(request)
	return fmt.Sprintf(callbackURLPattern, scheme, host, builder.serverBasePath, provider)
}

func (*HTTPThirdPartyCallbackURLBuilder) getRequestScheme(request *http.Request) string {
	if request.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}

func NewHTTPThirdPartyCallbackURLBuilder(serverBasePath string) *HTTPThirdPartyCallbackURLBuilder {
	return &HTTPThirdPartyCallbackURLBuilder{
		serverBasePath: serverBasePath,
	}
}
