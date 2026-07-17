package fake

import (
	"bytes"
	"io"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type FileStore struct {
	GotArtefact      bool
	ErrOnStore       error
	ErrOnGet         error
	NumDeletedImages int
	Data             []byte
	GotData          []byte
}

func (r *FileStore) Store(id im.ImageId, reader io.Reader) error {
	if r.ErrOnStore != nil {
		return r.ErrOnStore
	}
	r.GotArtefact = true
	data, _ := io.ReadAll(reader)
	r.GotData = data
	return nil
}

func (r *FileStore) Delete(im.ImageId) error {
	r.NumDeletedImages += 1
	return nil
}

func (r *FileStore) Get(im.ImageId) (io.Reader, error) {
	if r.ErrOnGet != nil {
		return nil, r.ErrOnGet
	}
	return bytes.NewBuffer(r.Data), nil
}
