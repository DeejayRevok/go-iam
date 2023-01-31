package app

import (
	"fmt"
	"os"

	"github.com/mvrilo/go-redoc"
)

func NewRedocConfiguration() *redoc.Redoc {
	basePath := os.Getenv("IAM_HTTP_SERVER_BASE_PATH")
	doc := redoc.Redoc{
		Title:       "iam",
		Description: "iam OpenAPI spec",
		SpecFile:    "./app/openapiDefinitions/open_api_spec_0.1.0.yaml",
		SpecPath:    fmt.Sprintf("%s/openapi.yaml", basePath),
		DocsPath:    "/docs",
	}
	return &doc
}
