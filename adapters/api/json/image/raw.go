package image

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

type Raw struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Raw) SuccessReadRawImage(r raw.Response) {

	data, err := io.ReadAll(r.Reader)
	if err != nil {
		http.Error(p.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	sum := sha256.Sum256(data)
	etag := `"` + hex.EncodeToString(sum[:]) + `"`
	p.Writer.Header().Set("ETag", etag)
	p.Writer.Header().Set("Content-Type", r.MIMEType)
	p.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	p.Writer.Header().Set("Cache-Control", "private, max-age=3600")
	p.Writer.WriteHeader(http.StatusOK)
	p.Writer.Write(data)
}

func NewRawImagePresenter(w http.ResponseWriter, l slog.Logger) Raw {
	return Raw{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
