package cli

import (
	"log/slog"
	"os"
)

type ErrorPresenter struct {
	slog.Logger
}

func NewErrorPresenter() ErrorPresenter {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return ErrorPresenter{*logger}
}

func (p ErrorPresenter) Error(err error) {
	p.Logger.Error(err.Error())
}
