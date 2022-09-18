package logging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type ZapGormTracedLogger struct {
	zapTracedLogger *ZapTracedLogger
}

func (l *ZapGormTracedLogger) LogMode(level logger.LogLevel) logger.Interface {
	return nil
}

func (l *ZapGormTracedLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.zapTracedLogger.Info(ctx, msg, data)
}

func (l *ZapGormTracedLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.zapTracedLogger.Warn(ctx, msg, data)
}

func (l *ZapGormTracedLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.zapTracedLogger.Error(ctx, msg, data)
}

func (l *ZapGormTracedLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		sql, rows := fc()
		if rows == -1 {
			l.zapTracedLogger.Error(ctx, fmt.Sprintf("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			l.zapTracedLogger.Error(ctx, fmt.Sprintf("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	}
}

func NewZapGormTracedLogger(tracer *apm.Tracer, opts ...zap.Option) *ZapGormTracedLogger {
	opts = append(opts, zap.AddCallerSkip(1))
	zapTracedLogger := NewZapTracedLogger(tracer, opts...)
	return &ZapGormTracedLogger{
		zapTracedLogger: zapTracedLogger,
	}
}
