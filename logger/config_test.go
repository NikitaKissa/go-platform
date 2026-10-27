package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWithConfigOutputModes(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "app.log")
	tests := []struct {
		name   string
		config Config
	}{
		{name: "stdout", config: Config{Output: OutputStdout}},
		{name: "file", config: Config{Output: OutputFile, FilePath: filePath}},
		{name: "file and stdout", config: Config{Output: OutputFileAndStdout, FilePath: filePath}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, err := NewWithConfig(test.config)
			if err != nil {
				t.Fatalf("NewWithConfig() error = %v", err)
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

func TestNewWithConfigRequiresFilePath(t *testing.T) {
	for _, output := range []OutputMode{OutputFile, OutputFileAndStdout} {
		if _, err := NewWithConfig(Config{Output: output}); err == nil {
			t.Errorf("NewWithConfig(Output: %d) expected an error without a file path", output)
		}
	}
}
