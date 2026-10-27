package logger

import (
	"time"

	"go.uber.org/zap"
)

type zapLogger struct {
	log *zap.Logger
}

func New() Logger {
	logger, err := NewWithConfig(Config{})
	if err != nil {
		panic(err)
	}
	return logger
}

func NewWithConfig(config Config) (Logger, error) {
	outputPaths, err := config.outputPaths()
	if err != nil {
		return nil, err
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.OutputPaths = outputPaths
	log, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return &zapLogger{log: log}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.log.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.log.Info(msg, toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.log.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.log.Error(msg, toZapFields(fields)...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{log: l.log.With(toZapFields(fields)...)}
}

func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{log: l.log.Named(name)}
}

func (l *zapLogger) Sync() error {
	return l.log.Sync()
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for index, field := range fields {
		switch field.kind {
		case fieldString:
			zapFields[index] = zap.String(field.key, field.value.(string))
		case fieldInt:
			zapFields[index] = zap.Int(field.key, field.value.(int))
		case fieldInt64:
			zapFields[index] = zap.Int64(field.key, field.value.(int64))
		case fieldUint:
			zapFields[index] = zap.Uint(field.key, field.value.(uint))
		case fieldBool:
			zapFields[index] = zap.Bool(field.key, field.value.(bool))
		case fieldFloat64:
			zapFields[index] = zap.Float64(field.key, field.value.(float64))
		case fieldDuration:
			zapFields[index] = zap.Duration(field.key, field.value.(time.Duration))
		case fieldTime:
			zapFields[index] = zap.Time(field.key, field.value.(time.Time))
		case fieldError:
			zapFields[index] = zap.Error(field.value.(error))
		case fieldAny:
			zapFields[index] = zap.Any(field.key, field.value)
		}
	}
	return zapFields
}
