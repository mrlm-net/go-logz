package logger

import (
	"fmt"
	"io"
	"os"
)

// ConsoleOutput creates an output function that writes to stdout/stderr
func ConsoleOutput() OutputFunc {
	return func(level LogLevel, message string) {
		if level <= Error {
			fmt.Fprintln(os.Stderr, message)
		} else {
			fmt.Fprintln(os.Stdout, message)
		}
	}
}

// FileOutput creates an output function that writes to a file
func FileOutput(filename string) (OutputFunc, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filename, err)
	}

	return func(level LogLevel, message string) {
		fmt.Fprintln(file, message)
	}, nil
}

// WriterOutput creates an output function that writes to any io.Writer
func WriterOutput(writer io.Writer) OutputFunc {
	return func(level LogLevel, message string) {
		fmt.Fprintln(writer, message)
	}
}

// SplitOutput creates separate outputs for different log levels (e.g., errors to stderr, info to stdout)
func SplitOutput(errorWriter, infoWriter io.Writer) OutputFunc {
	return func(level LogLevel, message string) {
		if level <= Error {
			fmt.Fprintln(errorWriter, message)
		} else {
			fmt.Fprintln(infoWriter, message)
		}
	}
}

// MultiOutput combines multiple outputs into one
func MultiOutput(outputs ...OutputFunc) OutputFunc {
	return func(level LogLevel, message string) {
		for _, output := range outputs {
			output(level, message)
		}
	}
}

// LevelFilterOutput wraps an output with level filtering
func LevelFilterOutput(minLevel LogLevel, output OutputFunc) OutputFunc {
	return func(level LogLevel, message string) {
		if level <= minLevel {
			output(level, message)
		}
	}
}
