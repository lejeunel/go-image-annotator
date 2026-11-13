package fake

import (
	"bytes"
	"io"
)

// FileStore keeps stored files in memory, keyed by name.
type FileStore struct {
	ErrOnStore  error
	ErrOnGet    error
	ErrOnDelete error
	ErrOnExists error
	// Missing makes Exists report names that were never stored as missing.
	Missing         bool
	NumDeletedItems int
	NumGets         int
	// Data is returned for names that were never stored.
	Data []byte
	// GotData holds the content of the last Store call.
	GotData []byte
	files   map[string][]byte
}

func (r *FileStore) content(name string) []byte {
	if data, ok := r.files[name]; ok {
		return data
	}
	return r.Data
}

func (r *FileStore) Store(name string, reader io.Reader) error {
	if r.ErrOnStore != nil {
		return r.ErrOnStore
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if r.files == nil {
		r.files = map[string][]byte{}
	}
	r.files[name] = data
	r.GotData = data
	return nil
}

func (r *FileStore) GetReaderAt(name string) (io.ReaderAt, int64, error) {
	data := r.content(name)
	return bytes.NewReader(data), int64(len(data)), nil
}

func (r *FileStore) Delete(name string) error {
	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	delete(r.files, name)
	r.NumDeletedItems += 1
	return nil
}

func (r *FileStore) Get(name string) (io.ReadCloser, error) {
	r.NumGets += 1
	if r.ErrOnGet != nil {
		return nil, r.ErrOnGet
	}
	return io.NopCloser(bytes.NewReader(r.content(name))), nil
}

func (r *FileStore) Exists(name string) (bool, error) {
	if r.ErrOnExists != nil {
		return false, r.ErrOnExists
	}
	if _, ok := r.files[name]; ok {
		return true, nil
	}
	return !r.Missing, nil
}
