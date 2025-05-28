package main

import (
	"fmt"
	"time"

	"github.com/mrlm-net/go-logz/pkg/logger"
)

func main() {
	// Advanced Example: Web application logging with different outputs for different components

	// Create file outputs for different components
	apiLogFile, err := logger.FileOutput("./logs/api.log")
	if err != nil {
		fmt.Printf("Failed to create API log file: %v\n", err)
		return
	}

	dbLogFile, err := logger.FileOutput("./logs/database.log")
	if err != nil {
		fmt.Printf("Failed to create DB log file: %v\n", err)
		return
	}

	errorLogFile, err := logger.FileOutput("./logs/errors.log")
	if err != nil {
		fmt.Printf("Failed to create error log file: %v\n", err)
		return
	}

	// API Logger - logs API requests and responses
	apiLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Info,
		Format: logger.JSONOutput,
		Outputs: []logger.OutputFunc{
			logger.MultiOutput(
				logger.ConsoleOutput(),
				apiLogFile,
				logger.LevelFilterOutput(logger.Error, errorLogFile),
			),
		},
		Prefix: "API",
	})

	// Database Logger - logs database operations
	dbLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Debug,
		Format: logger.StringOutput,
		Outputs: []logger.OutputFunc{
			logger.MultiOutput(
				dbLogFile,
				logger.LevelFilterOutput(logger.Error, errorLogFile),
			),
		},
		Prefix: "DB",
	})

	// Application Logger - general application logs
	appLogger := logger.NewLogger(logger.LogOptions{
		Level:  logger.Info,
		Format: logger.StringOutput,
		Outputs: []logger.OutputFunc{
			logger.ConsoleOutput(),
			logger.LevelFilterOutput(logger.Error, errorLogFile),
		},
		Prefix: "APP",
	})

	// Custom formatter for structured API logs
	structuredAPILogger := logger.NewLogger(logger.LogOptions{
		Level: logger.Debug,
		FormatCallback: func(level logger.LogLevel, message string, additionalInfo map[string]interface{}) string {
			timestamp := time.Now().UTC().Format(time.RFC3339)

			// Build a structured log entry
			logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)

			if method, ok := additionalInfo["method"]; ok {
				if path, ok := additionalInfo["path"]; ok {
					if status, ok := additionalInfo["status"]; ok {
						logEntry = fmt.Sprintf("[%s] [%s] %s %s -> %v | %s",
							timestamp, level.String(), method, path, status, message)
					}
				}
			}

			return logEntry
		},
		Outputs: []logger.OutputFunc{
			logger.ConsoleOutput(),
			apiLogFile,
		},
		Prefix: "API-STRUCTURED",
	})

	// Simulate application workflow
	appLogger.Info("Application starting up")

	// Simulate API requests
	apiLogger.Info("Incoming API request", map[string]interface{}{
		"method":    "GET",
		"path":      "/api/users",
		"ip":        "192.168.1.100",
		"userAgent": "Mozilla/5.0",
	})

	structuredAPILogger.Info("Request processed", map[string]interface{}{
		"method":       "GET",
		"path":         "/api/users",
		"status":       200,
		"responseTime": "45ms",
	})

	// Simulate database operations
	dbLogger.Debug("Executing SQL query", map[string]interface{}{
		"query":    "SELECT * FROM users WHERE active = ?",
		"params":   []interface{}{true},
		"duration": "12ms",
	})

	dbLogger.Info("Query executed successfully", map[string]interface{}{
		"rowsReturned": 25,
		"duration":     "12ms",
	})

	// Simulate error scenarios
	apiLogger.Error("Authentication failed", map[string]interface{}{
		"method":   "POST",
		"path":     "/api/login",
		"ip":       "192.168.1.100",
		"reason":   "invalid_credentials",
		"attempts": 3,
	})

	dbLogger.Error("Database connection failed", map[string]interface{}{
		"host":     "localhost:5432",
		"database": "myapp",
		"error":    "connection timeout",
	})

	appLogger.Info("Application shutting down gracefully")

	fmt.Println("\nCheck the ./logs/ directory for log files:")
	fmt.Println("- api.log: API-related logs")
	fmt.Println("- database.log: Database operation logs")
	fmt.Println("- errors.log: All error-level messages from all loggers")
}
