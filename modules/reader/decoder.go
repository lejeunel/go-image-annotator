package reader

import (
	"encoding/base64"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Base64ImageDecoder struct {
	Base64Data string
	data       []byte
	offset     int
}

func NewBase64ImageDecoder(base64Data string) *Base64ImageDecoder {
	return &Base64ImageDecoder{Base64Data: base64Data}
}

func (r *Base64ImageDecoder) Read(p []byte) (int, error) {
	errCtx := "decoding base64 data"

	if r.data == nil {
		bytesData, err := base64.StdEncoding.DecodeString(r.Base64Data)
		if err != nil {
			return 0, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrValidation)
		}

		r.data = bytesData

	}
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.offset:])
	r.offset += n

	return n, nil
}
