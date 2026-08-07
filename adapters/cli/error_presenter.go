package cli

import (
	"log/slog"

	l "github.com/lejeunel/go-image-annotator/shared/logging"
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
