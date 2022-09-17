package app

import (
	"fmt"
	"go-uaa/src/infrastructure/tracing"
	"os"

	"go.uber.org/zap"
)

func Start() {
	container := BuildDIContainer()
	httpServer := BuildHTTPServer(&container)
	RunEventConsumers(&container)

	err := container.Invoke(func(logger *zap.Logger, tracerConfig *tracing.JaegerTracerConfig) {
		defer tracerConfig.TracerCloser.Close()
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
