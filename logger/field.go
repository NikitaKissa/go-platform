package logger

import "time"

type Field struct {
	key   string
	kind  fieldKind
	value any
}

type fieldKind uint8

const (
	fieldString fieldKind = iota
	fieldInt
	fieldInt64
	fieldUint
	fieldBool
	fieldFloat64
	fieldDuration
	fieldTime
	fieldError
	fieldAny
)

func String(key, value string) Field {
	return Field{key: key, kind: fieldString, value: value}
}

func Int(key string, value int) Field {
	return Field{key: key, kind: fieldInt, value: value}
}

func Int64(key string, value int64) Field {
	return Field{key: key, kind: fieldInt64, value: value}
}

func Uint(key string, value uint) Field {
	return Field{key: key, kind: fieldUint, value: value}
}

func Bool(key string, value bool) Field {
	return Field{key: key, kind: fieldBool, value: value}
}

func Float64(key string, value float64) Field {
	return Field{key: key, kind: fieldFloat64, value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{key: key, kind: fieldDuration, value: value}
}

func Time(key string, value time.Time) Field {
	return Field{key: key, kind: fieldTime, value: value}
}

func Error(err error) Field {
	return Field{key: "error", kind: fieldError, value: err}
}

func Any(key string, value any) Field {
	return Field{key: key, kind: fieldAny, value: value}
}
