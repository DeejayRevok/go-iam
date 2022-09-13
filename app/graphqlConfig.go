package app

import (
	"context"
	"go-uaa/src/infrastructure/graph/resolvers"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const schemasFolder = "graphqlSchemas"

func BuildGraphQLHTTPHandler(resolver *resolvers.RootResolver, logger *zap.Logger) echo.HandlerFunc {
	_, currentFile, _, _ := runtime.Caller(0)
	schemasFolderPath := path.Join(filepath.Dir(currentFile), schemasFolder)
	schema, err := joinSchemasFromFolder(schemasFolderPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	parsedSchema, err := graphql.ParseSchema(schema, resolver)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return getEchoHandler(parsedSchema)

}

func joinSchemasFromFolder(folderPath string) (string, error) {
	schemas, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return "", err
	}
	joinedSchemas := make([]byte, 0)
	for _, schemaFile := range schemas {
		schemaFileName := schemaFile.Name()
		if !schemaFile.IsDir() && filepath.Ext(schemaFileName) == ".graphql" {
			schemaBytes, err := ioutil.ReadFile(path.Join(folderPath, schemaFileName))
			if err != nil {
				return "", err
			}
			joinedSchemas = append(joinedSchemas, schemaBytes...)
		}
	}
	return string(joinedSchemas[:]), nil
}

func getEchoHandler(schema *graphql.Schema) echo.HandlerFunc {
	handler := &relay.Handler{Schema: schema}
	return func(c echo.Context) error {
		request := c.Request()
		ctx := request.Context()
		ctx = context.WithValue(ctx, resolvers.RequestKey, request)
		handler.ServeHTTP(c.Response().Writer, c.Request().WithContext(ctx))
		return nil
	}
}
