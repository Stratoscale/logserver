// +build !appengine,go1.7

package ctxlog

import (
	"context"
	"log"
)

func LoggerFromContext(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(LoggerContextKey).(*log.Logger)
	if !ok {
		logger = Logger
	}
	return logger
}

func PrefixFromContext(ctx context.Context) string {
	prefix, _ := ctx.Value(PrefixContextKey).(string)
	return prefix
}

func NewContext(ctx context.Context, prefix string) context.Context {
	prefix = PrefixFromContext(ctx) + prefix
	ctx = context.WithValue(ctx, PrefixContextKey, prefix)
	return ctx
}

func logf(ctx context.Context, level string, format string, args ...interface{}) {
	logger := LoggerFromContext(ctx)
	prefix := PrefixFromContext(ctx)
	args = append([]interface{}{prefix}, args...)
	logger.Printf("%s"+level+": "+format, args...)
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
