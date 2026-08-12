package ingest

import (
	"io"
)

type MaxBytesReader struct {
	io.Reader
	MaxNBytes      int64
	read           int64
	excessBytesErr error
}

func NewMaxBytesReader(r io.Reader, maxNBytes int64, excessBytesErr error) MaxBytesReader {
	return MaxBytesReader{Reader: r, MaxNBytes: maxNBytes, excessBytesErr: excessBytesErr}
}

func (r *MaxBytesReader) Read(p []byte) (int, error) {
	if r.read >= r.MaxNBytes {
		return 0, r.excessBytesErr
	}

	if int64(len(p)) > r.MaxNBytes-r.read+1 {
		p = p[:r.MaxNBytes-r.read+1]
	}

	n, err := r.Reader.Read(p)
	r.read += int64(n)

	if r.read > r.MaxNBytes {
		return n, r.excessBytesErr
	}

	return n, err
}
