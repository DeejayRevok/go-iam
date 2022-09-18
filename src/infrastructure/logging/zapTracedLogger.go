package logging

import (
	"context"

	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

type ZapTracedLogger struct {
	zapLogger *zap.Logger
}

func (l *ZapTracedLogger) Info(ctx context.Context, msg string, extraFields ...interface{}) {
	traceContextFields := apmzap.TraceContext(ctx)
	l.zapLogger.With(traceContextFields...).Info(msg)
}

func (l *ZapTracedLogger) Warn(ctx context.Context, msg string, extraFields ...interface{}) {
	traceContextFields := apmzap.TraceContext(ctx)
	l.zapLogger.With(traceContextFields...).Warn(msg)
}

func (l *ZapTracedLogger) Error(ctx context.Context, msg string, extraFields ...interface{}) {
	traceContextFields := apmzap.TraceContext(ctx)
	l.zapLogger.With(traceContextFields...).Error(msg)
}

func NewZapTracedLogger(tracer *apm.Tracer, opts ...zap.Option) *ZapTracedLogger {
	opts = append(opts, zap.AddCallerSkip(1))
	zapLogger := NewZapLogger(tracer, opts...)
	return &ZapTracedLogger{
		zapLogger: zapLogger,
	}
}
