package fake

import (
	"io"
	"log/slog"
)

type Logger struct{}

func NewLogger() slog.Logger {
	return *slog.New(slog.NewTextHandler(io.Discard, nil))

}
