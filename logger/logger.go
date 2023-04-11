// Deprecated: Package logger provides a interface for logging.
package logger

import (
	"context"
)

type Logger interface {
	Debug(v ...interface{})
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnln(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})
}

type nopLogger struct {
}

func NewNop() Logger {
	return &nopLogger{}
}

func (l *nopLogger) Debug(v ...interface{})                 {}
func (l *nopLogger) Debugln(v ...interface{})               {}
func (l *nopLogger) Debugf(format string, v ...interface{}) {}

func (l *nopLogger) Info(v ...interface{})                 {}
func (l *nopLogger) Infoln(v ...interface{})               {}
func (l *nopLogger) Infof(format string, v ...interface{}) {}

func (l *nopLogger) Warn(v ...interface{})                 {}
func (l *nopLogger) Warnln(v ...interface{})               {}
func (l *nopLogger) Warnf(format string, v ...interface{}) {}

func (l *nopLogger) Error(v ...interface{})                 {}
func (l *nopLogger) Errorln(v ...interface{})               {}
func (l *nopLogger) Errorf(format string, v ...interface{}) {}

type contextKey string

const loggerContextKey = "logger"

// Deprecated: ContextWithLogger no longer does anything.
func ContextWithLogger(parent context.Context, l Logger) context.Context {
	return context.WithValue(parent, loggerContextKey, l)
}

// Deprecated: FromContext only returns nop logger.
func FromContext(ctx context.Context) Logger {
	return ctx.Value(loggerContextKey).(Logger)
}
