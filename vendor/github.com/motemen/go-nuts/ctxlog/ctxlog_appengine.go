// +build appengine

package ctxlog

import (
	"log"

	"golang.org/x/net/context"
	aelog "google.golang.org/appengine/log"
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

func Debugf(ctx context.Context, format string, args ...interface{}) {
	prefix := FromContext(ctx).Prefix()
	args = append([]interface{}{prefix}, args...)
	aelog.Debugf(ctx, "%s"+format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	prefix := FromContext(ctx).Prefix()
	args = append([]interface{}{prefix}, args...)
	aelog.Infof(ctx, "%s"+format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	prefix := FromContext(ctx).Prefix()
	args = append([]interface{}{prefix}, args...)
	aelog.Warningf(ctx, "%s"+format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	prefix := FromContext(ctx).Prefix()
	args = append([]interface{}{prefix}, args...)
	aelog.Errorf(ctx, "%s"+format, args...)
}

func Criticalf(ctx context.Context, format string, args ...interface{}) {
	prefix := FromContext(ctx).Prefix()
	args = append([]interface{}{prefix}, args...)
	aelog.Criticalf(ctx, "%s"+format, args...)
}
