package fake

import (
	"bytes"
	"io"
)

type FileStore struct {
	ErrOnStore       error
	ErrOnGet         error
	NumDeletedImages int
	Data             []byte
	GotData          []byte
}

func (r *FileStore) Store(path string, reader io.Reader) error {
	if r.ErrOnStore != nil {
		return r.ErrOnStore
	}
	data, _ := io.ReadAll(reader)
	r.GotData = data
	return nil
}
func (r *FileStore) GetReaderAt(string) (io.ReaderAt, int64, error) {
	return bytes.NewReader(r.Data), int64(len(r.Data)), nil
}

func (r *FileStore) Delete(string) error {
	r.NumDeletedImages += 1
	return nil
}

func (r *FileStore) Get(string) (io.Reader, error) {
	if r.ErrOnGet != nil {
		return nil, r.ErrOnGet
	}
	return bytes.NewBuffer(r.Data), nil
}
