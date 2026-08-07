package logging

import (
	"io"
	"log/slog"
	"os"
)

func NewNoOpLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func NewCliLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
