package logger

import (
	"fmt"
	"log/slog"

	"ai-search-bot/internal/config"
)

func New(env string, cfg config.LogConfig) (*slog.Logger, func(), error) {
	level := parseLevel(cfg.Level)

	switch env {
	case "dev":
		l, cleanup := newConsole(level)
		return l, cleanup, nil
	default:
		if cfg.File == "" {
			return nil, nil, fmt.Errorf("log.file must be set for env %q", env)
		}
		return newFile(cfg.File, level)
	}
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
