package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	logLevelDebug   = "DEBUG"
	logLevelInfo    = "INFO"
	logLevelWarn    = "WARN"
	logLevelWarning = "WARNING"
	logLevelError   = "ERROR"

	logFormatPlain = "PLAIN"
	logFormatJSON  = "JSON"
)

// Logger contains the service logger configuration properties.
type Logger struct {
	// Level is the log level (DEBUG, INFO, WARN, WARNING, ERROR).
	Level string `yaml:"level,omitempty" json:"level,omitempty"`
	// Format is the log format (PLAIN, JSON).
	Format string `yaml:"format,omitempty" json:"format,omitempty"`
}

// NewLoggerDefault returns a new Logger with default values.
func NewLoggerDefault() *Logger {
	return &Logger{
		Level:  logLevelInfo,
		Format: logFormatPlain,
	}
}

var (
	validLoggerLevels = []string{logLevelDebug, logLevelInfo, logLevelWarn,
		logLevelWarning, logLevelError}
	supportedLoggerFormats = []string{logFormatPlain, logFormatJSON}
)

func (l *Logger) SlogHandler() (slog.Handler, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}
	logLevel, err := l.logLevel()
	if err != nil {
		return nil, err
	}
	addSource := true
	writer := os.Stdout
	switch strings.ToUpper(l.Format) {
	case logFormatPlain:
		return slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: addSource,
		}), nil
	case logFormatJSON:
		return slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: addSource,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported log format: %s", l.Format)
	}
}

// validate validates the logger configuration properties.
func (l *Logger) validate() error {
	if l == nil {
		return errors.New("logger config is nil")
	}
	if !slices.Contains(validLoggerLevels, strings.ToUpper(l.Level)) {
		return fmt.Errorf("unsupported log level: %s", l.Level)
	}
	if !slices.Contains(supportedLoggerFormats, strings.ToUpper(l.Format)) {
		return fmt.Errorf("unsupported log format: %s", l.Format)
	}
	return nil
}

// logLevel returns the log level.
func (l *Logger) logLevel() (slog.Level, error) {
	switch strings.ToUpper(l.Level) {
	case logLevelDebug:
		return slog.LevelDebug, nil
	case logLevelInfo:
		return slog.LevelInfo, nil
	case logLevelWarn, logLevelWarning:
		return slog.LevelWarn, nil
	case logLevelError:
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", l.Level)
	}
}
