package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	logConfig := zap.NewProductionConfig()
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.TimeKey = "timestamp"
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig = logEncoderConfig

	log, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	return log
}
