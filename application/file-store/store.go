package file_store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

type FileStore struct {
	baseDir string
}

func NewFileStore(baseDir string) *FileStore {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create base directory: %v", err))
	}
	return &FileStore{baseDir: baseDir}
}

func (r *FileStore) filePath(id im.ImageId) string {
	return filepath.Join(r.baseDir, fmt.Sprintf("%s", id.String()))
}

func (r *FileStore) Store(id im.ImageId, reader io.Reader) error {
	path := r.filePath(id)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (r *FileStore) Delete(id im.ImageId) error {
	path := r.filePath(id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("artefact not found: %w", err)
		}
		return err
	}
	return nil
}

func (r *FileStore) Get(id im.ImageId) (io.Reader, error) {
	path := r.filePath(id)
	reader, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w: %w", err, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w", err, e.ErrInternal)
	}
	return reader, nil
}
