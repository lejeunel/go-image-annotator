package fake

import (
	"bytes"
)

type ImageReader struct {
	Buffer bytes.Buffer
	Err    error
}

func (d *ImageReader) Read(b []byte) (int, error) {
	if d.Err != nil {
		return 0, d.Err
	}
	return d.Buffer.Read(b)
}
