package logger

import (
	"fmt"
)

type OutputMode uint8

const (
	OutputStdout OutputMode = iota
	OutputFile
	OutputFileAndStdout
)

type Level uint8

const (
	LevelInfo Level = iota
	LevelDebug
	LevelWarn
	LevelError
)

type Encoding uint8

const (
	EncodingJSON Encoding = iota
	EncodingConsole
)

type Config struct {
	Output            OutputMode
	FilePath          string
	Level             Level
	Encoding          Encoding
	Development       bool
	AddCaller         bool
	DisableStacktrace bool
}

func (config Config) outputPaths() ([]string, error) {
	switch config.Output {
	case OutputStdout:
		return []string{"stdout"}, nil
	case OutputFile:
		if config.FilePath == "" {
			return nil, fmt.Errorf("logger: file path is required for file output")
		}
		return []string{config.FilePath}, nil
	case OutputFileAndStdout:
		if config.FilePath == "" {
			return nil, fmt.Errorf("logger: file path is required for file and stdout output")
		}
		return []string{"stdout", config.FilePath}, nil
	default:
		return nil, fmt.Errorf("logger: unsupported output mode %d", config.Output)
	}
}
