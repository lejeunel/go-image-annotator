package generic

import (
	e "datahub/errors"
	"log/slog"
	"net/http"
)

func LogAndWriteError(logger *slog.Logger, err error, w *http.ResponseWriter) {
	logger.Error(err.Error())
	http.Error(*w, err.Error(), e.ToHTTPStatus(err))
}
