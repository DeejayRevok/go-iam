package logging

import (
	"os"

	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(tracer *apm.Tracer, opts ...zap.Option) *zap.Logger {
	baseEncoderConfig := zap.NewProductionEncoderConfig()
	baseEncoderConfig.TimeKey = "timestamp"
	baseEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	stdCore := getStdCore(baseEncoderConfig)
	apmCore := getAPMCore(tracer)
	fileCore := getFileCore(baseEncoderConfig)
	logCore := zapcore.NewTee(stdCore, fileCore, apmCore)

	return zap.New(logCore, opts...)
}

func getAPMCore(tracer *apm.Tracer) zapcore.Core {
	return &apmzap.Core{
		Tracer: tracer,
	}
}

func getFileCore(encoderConfig zapcore.EncoderConfig) zapcore.Core {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	})
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriteSyncer,
		zap.InfoLevel,
	)
}

func getStdCore(encoderConfig zapcore.EncoderConfig) zapcore.Core {
	stdConfig := zap.NewProductionConfig()
	stdConfig.EncoderConfig = encoderConfig
	log, err := stdConfig.Build()
	if err != nil {
		panic(err)
	}
	return log.Core()
}
