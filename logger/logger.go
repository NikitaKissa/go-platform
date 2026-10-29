package logger

import (
	"context"
	"errors"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)

	With(fields ...Field) Logger
	Named(name string) Logger

	Sync() error
}

type loggerKey struct{}

var ErrContextNotFound = errors.New("logger not found in context")

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) (*Logger, error) {
	log, ok := ctx.Value(loggerKey{}).(*Logger)
	if !ok || log == nil {
		return nil, ErrContextNotFound
	}

	return log, nil
}
