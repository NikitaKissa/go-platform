package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.Logger
}

func New(config Config) (Logger, error) {
	outputPaths, err := config.outputPaths()
	if err != nil {
		return nil, err
	}

	level, err := config.Level.zapLevel()
	if err != nil {
		return nil, err
	}
	encoding, err := config.Encoding.zapEncoding()
	if err != nil {
		return nil, err
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       config.Development,
		DisableCaller:     !config.AddCaller,
		DisableStacktrace: config.DisableStacktrace,
		Encoding:          encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}
	log, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return &zapLogger{log: log}, nil
}

func (level Level) zapLevel() (zapcore.Level, error) {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel, nil
	case LevelDebug:
		return zapcore.DebugLevel, nil
	case LevelWarn:
		return zapcore.WarnLevel, nil
	case LevelError:
		return zapcore.ErrorLevel, nil
	default:
		return 0, fmt.Errorf("logger: unsupported level %d", level)
	}
}

func (encoding Encoding) zapEncoding() (string, error) {
	switch encoding {
	case EncodingJSON:
		return "json", nil
	case EncodingConsole:
		return "console", nil
	default:
		return "", fmt.Errorf("logger: unsupported encoding %d", encoding)
	}
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
