// +build !appengine,!go1.7

package ctxlog

import (
	"log"

	"golang.org/x/net/context"
)

func FromContext(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(LoggerContextKey).(*log.Logger)
	if !ok {
		logger = Logger
	}
	return logger
}

func NewContext(ctx context.Context, prefix string) context.Context {
	logger := FromContext(ctx)
	newLogger := log.New(output, logger.Prefix()+prefix, logger.Flags())
	return context.WithValue(ctx, LoggerContextKey, newLogger)
}

func logf(ctx context.Context, level string, format string, args ...interface{}) {
	logger := FromContext(ctx)
	args = append([]interface{}{level}, args...)
	logger.Printf("%s: "+format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logf(ctx, "debug", format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logf(ctx, "info", format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	logf(ctx, "warning", format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logf(ctx, "error", format, args...)
}

func Criticalf(ctx context.Context, format string, args ...interface{}) {
	logf(ctx, "critical", format, args...)
}
