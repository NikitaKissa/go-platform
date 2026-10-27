package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewOutputModes(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "app.log")
	tests := []struct {
		name   string
		config Config
	}{
		{name: "stdout", config: testConfig(OutputStdout, "")},
		{name: "file", config: testConfig(OutputFile, filePath)},
		{name: "file and stdout", config: testConfig(OutputFileAndStdout, filePath)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, err := New(test.config)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if err := logger.Sync(); err != nil && test.config.Output == OutputFile {
				t.Fatalf("Sync() error = %v", err)
			}
		})
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected file output to create %q: %v", filePath, err)
	}
}

func TestNewRequiresFilePath(t *testing.T) {
	for _, output := range []OutputMode{OutputFile, OutputFileAndStdout} {
		if _, err := New(testConfig(output, "")); err == nil {
			t.Errorf("New(Output: %d) expected an error without a file path", output)
		}
	}
}

func TestNewRejectsUnsupportedOptions(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "level", config: Config{Level: Level(255), Encoding: EncodingJSON}},
		{name: "encoding", config: Config{Level: LevelInfo, Encoding: Encoding(255)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := New(test.config); err == nil {
				t.Fatal("New() expected an error for unsupported option")
			}
		})
	}
}

func TestNewUsesProvidedLevelAndEncoding(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "custom.log")
	config := testConfig(OutputFile, filePath)
	config.Level = LevelWarn
	config.Encoding = EncodingJSON

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	logger.Info("filtered message")
	logger.Warn("written message")
	if err := logger.Sync(); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	contents, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := string(contents); got == "" || !contains(got, `"msg":"written message"`) || contains(got, "filtered message") {
		t.Fatalf("New() did not honor provided level and encoding; output = %q", got)
	}
}

func testConfig(output OutputMode, filePath string) Config {
	return Config{
		Output:            output,
		FilePath:          filePath,
		Level:             LevelDebug,
		Encoding:          EncodingJSON,
		DisableStacktrace: true,
	}
}

func contains(value, substring string) bool {
	return strings.Contains(value, substring)
}
