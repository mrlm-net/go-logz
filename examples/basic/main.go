package main

import (
	"os"

	"github.com/mrlm-net/go-logz/pkg/logger"
)

func main() {
	// Example 1: Basic console logging with string format
	basicLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.StringOutput,
		Prefix: "my-app",
	})

	basicLogger.Info("This is an info message")
	basicLogger.Error("This is an error message", map[string]interface{}{
		"errorCode": 123,
		"userId":    "user-456",
	})

	// Example 2: JSON format logging
	jsonLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.JSONOutput,
		Prefix: "my-app",
	})

	jsonLogger.Info("This is a JSON formatted message")
	jsonLogger.Warning("API rate limit approaching", map[string]interface{}{
		"currentRate": 95,
		"limit":       100,
		"endpoint":    "/api/users",
	})

	// Example 3: File outputs
	stdoutFile, _ := logger.FileOutput("./stdout.log")
	stderrFile, _ := logger.FileOutput("./stderr.log")

	fileLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.StringOutput,
		Outputs: []logger.OutputFunc{
			func(level logger.LogLevel, message string) {
				if level <= logger.Error {
					stderrFile(level, message)
				} else {
					stdoutFile(level, message)
				}
			},
		},
		Prefix: "my-app",
	})

	fileLogger.Info("This goes to stdout.log")
	fileLogger.Error("This goes to stderr.log", map[string]interface{}{
		"errorCode": 500,
	})

	// Example 4: Multiple outputs (console + file)
	allOutputFile, _ := logger.FileOutput("./all.log")
	multiLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.StringOutput,
		Outputs: []logger.OutputFunc{
			logger.MultiOutput(
				logger.ConsoleOutput(),
				allOutputFile,
			),
		},
		Prefix: "my-app",
	})

	multiLogger.Info("This message goes to both console and file")

	// Example 5: Custom format callback
	customLogger := logger.NewLogger(logger.LogOptions{
		Level: logger.Debug,
		FormatCallback: func(level logger.LogLevel, message string, additionalInfo map[string]interface{}) string {
			return "🚀 [CUSTOM] " + level.String() + ": " + message
		},
		Prefix: "custom",
	})

	customLogger.Info("Custom formatted message")

	// Example 6: Level filtering
	levelFilterLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.StringOutput,
		Outputs: []logger.OutputFunc{
			logger.LevelFilterOutput(logger.Error, logger.WriterOutput(os.Stderr)),
			logger.LevelFilterOutput(logger.Info, logger.WriterOutput(os.Stdout)),
		},
		Prefix: "filtered",
	})

	levelFilterLogger.Debug("This won't appear anywhere")
	levelFilterLogger.Info("This goes to stdout")
	levelFilterLogger.Error("This goes to stderr")
}
