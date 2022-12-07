package app

import (
	"fmt"
	"net/url"
	"os"

	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
)

func NewAPMTracer() *apm.Tracer {
	tracerTransport, err := getHTTPTransport()
	if err != nil {
		panic(fmt.Sprintf("Error initializing APM tracer transport: %s", err.Error()))
	}
	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:    "iam",
		ServiceVersion: "1",
		Transport:      tracerTransport,
	})
	if err != nil {
		panic(fmt.Sprintf("Error initializing APM tracer: %s", err.Error()))
	}
	apm.SetDefaultTracer(tracer)
	return tracer
}

func getHTTPTransport() (*transport.HTTPTransport, error) {
	transportOptions := getHTTPTransportOptions()

	httpTransport, err := transport.NewHTTPTransport(*transportOptions)
	if err != nil {
		return nil, err
	}
	return httpTransport, nil
}

func getHTTPTransportOptions() *transport.HTTPTransportOptions {
	apmServerHost := os.Getenv("ELASTIC_APM_HOST")
	apmServerPort := os.Getenv("ELASTIC_APM_PORT")
	apmSecretToken := os.Getenv("ELASTIC_APM_SECRET_TOKEN")

	apmServerURL := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", apmServerHost, apmServerPort),
	}

	return &transport.HTTPTransportOptions{
		ServerURLs:  []*url.URL{&apmServerURL},
		SecretToken: apmSecretToken,
	}
}
