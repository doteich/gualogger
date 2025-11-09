package logging

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

// KgoLoggerAdapter adapts a *slog.Logger to the kgo.Logger interface.
type KgoLoggerAdapter struct {
	slog *slog.Logger
}

// NewKgoLogger is a constructor for the adapter.
func NewKgoLogger(logger *slog.Logger) *KgoLoggerAdapter {
	return &KgoLoggerAdapter{slog: logger}
}

// Level determines the kgo.LogLevel from the underlying slog.Logger's enabled level.
func (a *KgoLoggerAdapter) Level() kgo.LogLevel {
	if a.slog.Enabled(context.Background(), slog.LevelDebug) {
		return kgo.LogLevelDebug
	}
	if a.slog.Enabled(context.Background(), slog.LevelInfo) {
		return kgo.LogLevelInfo
	}
	if a.slog.Enabled(context.Background(), slog.LevelWarn) {
		return kgo.LogLevelWarn
	}
	if a.slog.Enabled(context.Background(), slog.LevelError) {
		return kgo.LogLevelError
	}
	return kgo.LogLevelNone
}

// Log forwards the log message from kgo to the underlying slog.Logger.
func (a *KgoLoggerAdapter) Log(level kgo.LogLevel, msg string, keyvals ...interface{}) {
	var slogLevel slog.Level
	switch level {
	case kgo.LogLevelError:
		slogLevel = slog.LevelError
	case kgo.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case kgo.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case kgo.LogLevelDebug:
		slogLevel = slog.LevelDebug
	default:
		// Do not log for kgo.LogLevelNone
		return
	}
	a.slog.Log(context.Background(), slogLevel, msg, keyvals...)
}
