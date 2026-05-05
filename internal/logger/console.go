package logger

import (
	"log/slog"
	"os"
)

func newConsole(level slog.Level) (*slog.Logger, func()) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler), func() {}
}
