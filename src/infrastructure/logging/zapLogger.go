package logging

import (
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(tracer *apm.Tracer, opts ...zap.Option) *zap.Logger {
	apmZapCore := apmzap.Core{
		Tracer: tracer,
	}
	zapCore := zap.WrapCore((&apmZapCore).WrapCore)
	logConfig := zap.NewProductionConfig()
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.TimeKey = "timestamp"
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig = logEncoderConfig

	opts = append(opts, zapCore)
	log, err := logConfig.Build(opts...)
	if err != nil {
		panic(err)
	}
	return log
}
