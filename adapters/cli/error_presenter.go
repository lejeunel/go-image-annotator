package cli

import (
	l "github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type ErrorPresenter struct {
	slog.Logger
}

func NewErrorPresenter() ErrorPresenter {
	return ErrorPresenter{*l.NewCliLogger()}
}

func (p ErrorPresenter) Error(err error) {
	p.Logger.Error(err.Error())
}
