// Package logger provides SysLog RFC 5424 compliant logging functionality
// with support for multiple output formats and destinations.
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogLevel represents the severity level of log messages following RFC 5424
type LogLevel int

const (
	Emergency LogLevel = iota // System is unusable
	Alert                     // Action must be taken immediately
	Critical                  // Critical conditions
	Error                     // Error conditions
	Warning                   // Warning conditions
	Notice                    // Normal but significant condition
	Info                      // Informational messages
	Debug                     // Debug-level messages
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case Emergency:
		return "EMERGENCY"
	case Alert:
		return "ALERT"
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Notice:
		return "NOTICE"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// LogOutput represents the format of log messages
type LogOutput string

const (
	StringOutput LogOutput = "string"
	JSONOutput   LogOutput = "json"
)

// OutputFunc is a function that handles where logs are written
type OutputFunc func(level LogLevel, message string)

// FormatFunc is a custom formatting function
type FormatFunc func(level LogLevel, message string, additionalInfo map[string]interface{}) string

// LogOptions contains configuration options for the logger
type LogOptions struct {
	Level          LogLevel     // Minimum log level to output
	Format         LogOutput    // Output format (string or json)
	FormatCallback FormatFunc   // Custom format function
	Outputs        []OutputFunc // Output destination functions
	Prefix         string       // Prefix for log messages
}

// ILogger interface defines the logging contract
type ILogger interface {
	Log(level LogLevel, message string, additionalInfo ...map[string]interface{})
	Emergency(message string, additionalInfo ...map[string]interface{})
	Alert(message string, additionalInfo ...map[string]interface{})
	Critical(message string, additionalInfo ...map[string]interface{})
	Error(message string, additionalInfo ...map[string]interface{})
	Warning(message string, additionalInfo ...map[string]interface{})
	Notice(message string, additionalInfo ...map[string]interface{})
	Info(message string, additionalInfo ...map[string]interface{})
	Debug(message string, additionalInfo ...map[string]interface{})
}

// Logger implements the ILogger interface
type Logger struct {
	level          LogLevel
	format         LogOutput
	formatCallback FormatFunc
	outputs        []OutputFunc
	prefix         string
}

// NewLogger creates a new logger instance with the given options
func NewLogger(options LogOptions) *Logger {
	logger := &Logger{
		level:          options.Level,
		format:         options.Format,
		formatCallback: options.FormatCallback,
		outputs:        options.Outputs,
		prefix:         options.Prefix,
	}

	// Set defaults
	if logger.level < Emergency || logger.level > Debug {
		logger.level = Info
	}
	if logger.format == "" {
		logger.format = StringOutput
	}
	if len(logger.outputs) == 0 {
		logger.outputs = []OutputFunc{defaultConsoleOutput}
	}

	return logger
}

// defaultConsoleOutput is the default output function that writes to stdout/stderr
func defaultConsoleOutput(level LogLevel, message string) {
	if level <= Error {
		fmt.Fprint(os.Stderr, message+"\n")
	} else {
		fmt.Fprint(os.Stdout, message+"\n")
	}
}

// formatMessage formats the log message according to the configured format
func (l *Logger) formatMessage(level LogLevel, message string, additionalInfo map[string]interface{}) string {
	if l.formatCallback != nil {
		return l.formatCallback(level, message, additionalInfo)
	}

	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	logLevel := level.String()

	if l.format == JSONOutput {
		logData := map[string]interface{}{
			"timestamp": timestamp,
			"level":     logLevel,
			"message":   message,
		}
		if l.prefix != "" {
			logData["prefix"] = l.prefix
		}
		for k, v := range additionalInfo {
			logData[k] = v
		}

		jsonBytes, _ := json.Marshal(logData)
		return string(jsonBytes)
	}

	// String format
	var baseMessage string
	if l.prefix != "" {
		baseMessage = fmt.Sprintf("[%s] [%s] [%s]", l.prefix, timestamp, logLevel)
	} else {
		baseMessage = fmt.Sprintf("[%s] [%s]", timestamp, logLevel)
	}

	if len(additionalInfo) > 0 {
		additionalJSON, _ := json.Marshal(additionalInfo)
		baseMessage += fmt.Sprintf(" %s", string(additionalJSON))
	}

	return fmt.Sprintf("%s %s", baseMessage, message)
}

// shouldLog checks if the message should be logged based on the configured level
func (l *Logger) shouldLog(level LogLevel) bool {
	return level <= l.level
}

// outputMessage sends the message to all configured outputs
func (l *Logger) outputMessage(level LogLevel, message string) {
	for _, output := range l.outputs {
		output(level, message)
	}
}

// Log logs a message at the specified level
func (l *Logger) Log(level LogLevel, message string, additionalInfo ...map[string]interface{}) {
	if l.shouldLog(level) {
		var info map[string]interface{}
		if len(additionalInfo) > 0 {
			info = additionalInfo[0]
		} else {
			info = make(map[string]interface{})
		}
		formattedMessage := l.formatMessage(level, message, info)
		l.outputMessage(level, formattedMessage)
	}
}

// Emergency logs an emergency level message
func (l *Logger) Emergency(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Emergency, message, additionalInfo...)
}

// Alert logs an alert level message
func (l *Logger) Alert(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Alert, message, additionalInfo...)
}

// Critical logs a critical level message
func (l *Logger) Critical(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Critical, message, additionalInfo...)
}

// Error logs an error level message
func (l *Logger) Error(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Error, message, additionalInfo...)
}

// Warning logs a warning level message
func (l *Logger) Warning(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Warning, message, additionalInfo...)
}

// Notice logs a notice level message
func (l *Logger) Notice(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Notice, message, additionalInfo...)
}

// Info logs an info level message
func (l *Logger) Info(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Info, message, additionalInfo...)
}

// Debug logs a debug level message
func (l *Logger) Debug(message string, additionalInfo ...map[string]interface{}) {
	l.Log(Debug, message, additionalInfo...)
}
