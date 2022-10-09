package api

import (
	"fmt"
	"net/http"
)

const callbackURLPattern = "%s://%s/auth/%s/callback"

type HTTPThirdPartyCallbackURLBuilder struct {
}

func (builder *HTTPThirdPartyCallbackURLBuilder) Build(provider string, request *http.Request) string {
	host := request.Host
	scheme := builder.getRequestScheme(request)
	return fmt.Sprintf(callbackURLPattern, scheme, host, provider)
}

func (*HTTPThirdPartyCallbackURLBuilder) getRequestScheme(request *http.Request) string {
	if request.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}

func NewHTTPThirdPartyCallbackURLBuilder() *HTTPThirdPartyCallbackURLBuilder {
	return &HTTPThirdPartyCallbackURLBuilder{}
}
