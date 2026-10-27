package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var _ Logger = (*zapLogger)(nil)

func TestZapLoggerMethodsWriteMessagesAndFields(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "logger.log")
	logger, err := NewWithConfig(Config{Output: OutputFile, FilePath: filePath})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	if got := logger.Named("component").(*zapLogger).log.Name(); got != "component" {
		t.Fatalf("Named() logger name = %q; want %q", got, "component")
	}

	derived := logger.With(String("scope", "test")).Named("component")
	derived.Debug("debug message", String("level_field", "debug"))
	derived.Info("info message", Int("level_field", 1))
	derived.Warn("warn message", Bool("level_field", true))
	derived.Error("error message", Error(os.ErrNotExist))

	if err := derived.Sync(); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	contents, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	logOutput := string(contents)
	for _, expected := range []string{
		"component",
		"debug message",
		"info message",
		"warn message",
		"error message",
		`"scope": "test"`,
		`"level_field": "debug"`,
		`"level_field": 1`,
		`"level_field": true`,
		`"error": "file does not exist"`,
	} {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("log output does not contain %q:\n%s", expected, logOutput)
		}
	}
}

func TestNewReturnsLogger(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}
