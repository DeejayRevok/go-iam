package app

import (
	"github.com/mvrilo/go-redoc"
)

func NewRedocConfiguration() *redoc.Redoc {
	doc := redoc.Redoc{
		Title:       "UAA",
		Description: "UAA OpenAPI spec",
		SpecFile:    "./app/openapiDefinitions/open_api_spec_0.1.0.yaml",
		SpecPath:    "/openapi.yaml",
		DocsPath:    "/docs",
	}
	return &doc
}
