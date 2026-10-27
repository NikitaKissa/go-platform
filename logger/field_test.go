package logger

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestFieldConstructorsAndZapConversion(t *testing.T) {
	timestamp := time.Date(2026, time.October, 27, 12, 30, 0, 0, time.UTC)
	testError := errors.New("test error")
	tests := []struct {
		name    string
		field   Field
		kind    fieldKind
		key     string
		value   any
		zapType zapcore.FieldType
	}{
		{name: "string", field: String("key", "value"), kind: fieldString, key: "key", value: "value", zapType: zapcore.StringType},
		{name: "int", field: Int("key", 42), kind: fieldInt, key: "key", value: 42, zapType: zapcore.Int64Type},
		{name: "int64", field: Int64("key", 42), kind: fieldInt64, key: "key", value: int64(42), zapType: zapcore.Int64Type},
		{name: "uint", field: Uint("key", 42), kind: fieldUint, key: "key", value: uint(42), zapType: zapcore.Uint64Type},
		{name: "bool", field: Bool("key", true), kind: fieldBool, key: "key", value: true, zapType: zapcore.BoolType},
		{name: "float64", field: Float64("key", 3.5), kind: fieldFloat64, key: "key", value: 3.5, zapType: zapcore.Float64Type},
		{name: "duration", field: Duration("key", time.Second), kind: fieldDuration, key: "key", value: time.Second, zapType: zapcore.DurationType},
		{name: "time", field: Time("key", timestamp), kind: fieldTime, key: "key", value: timestamp, zapType: zapcore.TimeType},
		{name: "error", field: Error(testError), kind: fieldError, key: "error", value: testError, zapType: zapcore.ErrorType},
		{name: "any", field: Any("key", map[string]int{"count": 2}), kind: fieldAny, key: "key", value: map[string]int{"count": 2}, zapType: zapcore.ReflectType},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.field.kind != test.kind || test.field.key != test.key || !reflect.DeepEqual(test.field.value, test.value) {
				t.Fatalf("constructor returned %#v; want kind=%d key=%q value=%#v", test.field, test.kind, test.key, test.value)
			}

			converted := toZapFields([]Field{test.field})
			if len(converted) != 1 {
				t.Fatalf("toZapFields() returned %d fields; want 1", len(converted))
			}
			if converted[0].Key != test.key {
				t.Errorf("converted field key = %q; want %q", converted[0].Key, test.key)
			}
			if converted[0].Type != test.zapType {
				t.Errorf("converted field type = %v; want %v", converted[0].Type, test.zapType)
			}
		})
	}
}
