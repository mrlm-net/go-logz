package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerBasicFunctionality(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LogOptions{
		Level:  Debug,
		Format: StringOutput,
		Outputs: []OutputFunc{
			func(level LogLevel, message string) {
				buf.WriteString(message + "\n")
			},
		},
		Prefix: "test",
	})

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "[test]") {
		t.Errorf("Expected output to contain prefix '[test]', got: %s", output)
	}
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain '[INFO]', got: %s", output)
	}
}

func TestLoggerJSONFormat(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LogOptions{
		Level:  Debug,
		Format: JSONOutput,
		Outputs: []OutputFunc{
			func(level LogLevel, message string) {
				buf.WriteString(message + "\n")
			},
		},
		Prefix: "test",
	})

	logger.Info("test message", map[string]interface{}{
		"key": "value",
	})

	output := buf.String()
	if !strings.Contains(output, `"message":"test message"`) {
		t.Errorf("Expected JSON output to contain message field, got: %s", output)
	}
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected JSON output to contain level field, got: %s", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected JSON output to contain additional info, got: %s", output)
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LogOptions{
		Level:  Error, // Only log Error and above
		Format: StringOutput,
		Outputs: []OutputFunc{
			func(level LogLevel, message string) {
				buf.WriteString(message + "\n")
			},
		},
	})

	logger.Debug("debug message")         // Should not appear
	logger.Info("info message")           // Should not appear
	logger.Error("error message")         // Should appear
	logger.Emergency("emergency message") // Should appear

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Errorf("Debug message should not appear with Error level filter")
	}
	if strings.Contains(output, "info message") {
		t.Errorf("Info message should not appear with Error level filter")
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("Error message should appear with Error level filter")
	}
	if !strings.Contains(output, "emergency message") {
		t.Errorf("Emergency message should appear with Error level filter")
	}
}

func TestCustomFormatCallback(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LogOptions{
		Level: Debug,
		FormatCallback: func(level LogLevel, message string, additionalInfo map[string]interface{}) string {
			return "CUSTOM: " + message
		},
		Outputs: []OutputFunc{
			func(level LogLevel, message string) {
				buf.WriteString(message + "\n")
			},
		},
	})

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "CUSTOM: test message") {
		t.Errorf("Expected custom format, got: %s", output)
	}
}
