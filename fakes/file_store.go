package fake

import (
	"bytes"
	"io"
)

type FileStore struct {
	GotArtefact      bool
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
	r.GotArtefact = true
	data, _ := io.ReadAll(reader)
	r.GotData = data
	return nil
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
