# mrlm-net/go-logz

Logging package following SysLog protocol - [RFC 5424](https://datatracker.ietf.org/doc/html/rfc5424) written in Go.

| Package | `mrlm-net/go-logz` |
| :-- | :-- |
| Go Module | `github.com/mrlm-net/go-logz` |
| License | Apache 2.0 |

This is the Go implementation of the [TypeScript logz package](https://github.com/mrlm-net/logz).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Advanced Usage](#advanced-usage)
- [Output Destinations](#output-destinations)
- [Types](#types)
- [Examples](#examples)
- [Contributing](#contributing)

## Installation

```shell
go get github.com/mrlm-net/go-logz
```

## Usage

### Basic Usage

```go
package main

import (
    "github.com/mrlm-net/go-logz/pkg/logger"
)

func main() {
    // Create a basic logger
    log := logger.NewLogger(logger.LogOptions{
        Level:  logger.Debug,
        Format: logger.StringOutput,
        Prefix: "my-app",
    })

    log.Info("This is an info message")
    log.Error("This is an error message", map[string]interface{}{
        "errorCode": 123,
        "userId":    "user-456",
    })
}
```

Output:
```
[my-app] [2023-10-05T14:48:00.000Z] [INFO] This is an info message
[my-app] [2023-10-05T14:48:00.000Z] [ERROR] {"errorCode":123,"userId":"user-456"} This is an error message
```

### JSON Format

```go
jsonLogger := logger.NewLogger(logger.LogOptions{
    Level:  logger.Debug,
    Format: logger.JSONOutput,
    Prefix: "my-app",
})

jsonLogger.Info("This is a JSON message", map[string]interface{}{
    "requestId": "req-123",
    "duration":  "45ms",
})
```

Output:
```json
{"level":"INFO","message":"This is a JSON message","prefix":"my-app","requestId":"req-123","duration":"45ms","timestamp":"2023-10-05T14:48:00.000Z"}
```

## Advanced Usage

### Multiple Output Destinations

You can configure the logger to output logs to different destinations such as files, console, or any `io.Writer`:

```go
package main

import (
    "os"
    "github.com/mrlm-net/go-logz/pkg/logger"
)

func main() {
    // Create file outputs
    stdoutFile, _ := logger.FileOutput("./stdout.log")
    stderrFile, _ := logger.FileOutput("./stderr.log")

    // Logger that splits outputs based on log level
    log := logger.NewLogger(logger.LogOptions{
        Level:  logger.Debug,
        Format: logger.StringOutput,
        Outputs: []logger.OutputFunc{
            logger.SplitOutput(stderrFile, stdoutFile), // errors to stderr.log, others to stdout.log
        },
        Prefix: "my-app",
    })

    log.Info("This goes to stdout.log")
    log.Error("This goes to stderr.log")
}
```

### Multiple Outputs

```go
// Log to both console and file
multiLogger := logger.NewLogger(logger.LogOptions{
    Level:  logger.Debug,
    Format: logger.StringOutput,
    Outputs: []logger.OutputFunc{
        logger.MultiOutput(
            logger.ConsoleOutput(),
            stdoutFile,
        ),
    },
    Prefix: "my-app",
})

multiLogger.Info("This message goes to both console and file")
```

### Custom Format Callback

```go
customLogger := logger.NewLogger(logger.LogOptions{
    Level: logger.Debug,
    FormatCallback: func(level logger.LogLevel, message string, additionalInfo map[string]interface{}) string {
        return fmt.Sprintf("🚀 [CUSTOM] %s: %s", level.String(), message)
    },
    Prefix: "custom",
})

customLogger.Info("Custom formatted message")
// Output: 🚀 [CUSTOM] INFO: Custom formatted message
```

## Output Destinations

The package provides several built-in output functions:

### ConsoleOutput
Writes to stdout/stderr based on log level:
```go
logger.ConsoleOutput()
```

### FileOutput
Writes to a specified file:
```go
fileOutput, err := logger.FileOutput("./app.log")
if err != nil {
    // handle error
}
```

### WriterOutput
Writes to any `io.Writer`:
```go
logger.WriterOutput(os.Stdout)
```

### SplitOutput
Separates error and info logs to different writers:
```go
logger.SplitOutput(os.Stderr, os.Stdout)
```

### MultiOutput
Combines multiple outputs:
```go
logger.MultiOutput(
    logger.ConsoleOutput(),
    fileOutput,
)
```

### LevelFilterOutput
Filters output based on minimum log level:
```go
logger.LevelFilterOutput(logger.Error, logger.ConsoleOutput())
```

## Types

### LogLevel

```go
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
```

### LogOutput

```go
const (
    StringOutput LogOutput = "string"
    JSONOutput   LogOutput = "json"
)
```

### LogOptions

```go
type LogOptions struct {
    Level          LogLevel       // Minimum log level to output
    Format         LogOutput      // Output format (string or json)
    FormatCallback FormatFunc     // Custom format function
    Outputs        []OutputFunc   // Output destination functions
    Prefix         string         // Prefix for log messages
}
```

### ILogger Interface

```go
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
```

## Examples

See the [examples](./examples/) directory for complete working examples:

- [Basic Usage](./examples/basic/main.go) - Simple console and file logging
- [Advanced Usage](./examples/advanced/main.go) - Complex multi-output scenarios with structured logging

To run the examples:

```shell
# Basic example
go run ./examples/basic/main.go

# Advanced example (creates ./logs/ directory with log files)
mkdir -p logs
go run ./examples/advanced/main.go
```

## Key Features

- **RFC 5424 Compliant**: Follows SysLog severity levels
- **Multiple Output Formats**: String and JSON formatting
- **Flexible Output Destinations**: Console, files, custom writers
- **Structured Logging**: Support for additional context information
- **Custom Formatting**: Override default formatting with custom callbacks
- **Level Filtering**: Control which messages are logged based on severity
- **Thread-Safe**: Safe for concurrent use
- **Zero Dependencies**: Only uses Go standard library

## Contributing

_Contributions are welcomed and must follow [Code of Conduct](https://github.com/mrlm-net/logz?tab=coc-ov-file) and common [Contributions guidelines](https://github.com/mrlm-net/.github/blob/main/docs/CONTRIBUTING.md)._

> If you'd like to report security issue please follow [security guidelines](https://github.com/mrlm-net/logz?tab=security-ov-file).

---
<sup><sub>_All rights reserved © Martin Hrášek [<@marley-ma>](https://github.com/marley-ma) and WANTED.solutions s.r.o. [<@wanted-solutions>](https://github.com/wanted-solutions)_</sub></sup>
