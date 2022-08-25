package app

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func Start() {
	container := BuildDIContainer()
	httpServer := BuildHTTPServer(&container)
	RunEventConsumers(&container)

	err := container.Invoke(func(logger *zap.Logger) {
		serverHost := os.Getenv("HTTP_SERVER_HOST")
		serverPort := os.Getenv("HTTP_SERVER_PORT")
		handleError(httpServer.Start(fmt.Sprintf("%s:%s", serverHost, serverPort)), logger)
	})
	if err != nil {
		panic(err)
	}
}

func handleError(err error, logger *zap.Logger) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}
