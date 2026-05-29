package file_store

import (
	"bytes"
	"io"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type FakeStore struct {
	GotArtefact      bool
	Err              error
	NumDeletedImages int
	Data             []byte
	GotData          []byte
}

func (r *FakeStore) Store(id im.ImageId, reader io.Reader) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotArtefact = true
	data, _ := io.ReadAll(reader)
	r.GotData = data
	return nil
}

func (r *FakeStore) Delete(im.ImageId) error {
	r.NumDeletedImages += 1
	return nil
}

func (r *FakeStore) Get(im.ImageId) (io.Reader, error) {
	return bytes.NewBuffer(r.Data), nil
}
