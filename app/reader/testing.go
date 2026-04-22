package reader

import (
	_ "embed"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

//go:embed sample-image.jpg
var testJPGImage []byte

//go:embed sample-image.png
var testPNGImage []byte

type FakeFileStore struct {
	Err  error
	Data []byte
}

func (r *FakeFileStore) Get(im.ImageId) ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Data, nil
}
