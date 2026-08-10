package file_store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type LocalFileStore struct {
	baseDir string
}

func NewLocalFileStore(baseDir string) LocalFileStore {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		panic(fmt.Sprintf("failed to create base directory: %v", err))
	}
	return LocalFileStore{baseDir: baseDir}
}

func (r LocalFileStore) filePath(path string) string {
	return filepath.Join(r.baseDir, fmt.Sprintf("%s", path))
}

func (r LocalFileStore) Store(path string, reader io.Reader) error {
	path = r.filePath(path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

func (r LocalFileStore) Delete(path string) error {
	path = r.filePath(path)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("artefact not found: %w", err)
		}
		return err
	}
	return nil
}

func (r LocalFileStore) Get(path string) (io.Reader, error) {
	path = r.filePath(path)
	reader, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w: %w", err, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w", err, e.ErrInternal)
	}
	return reader, nil
}

func (r LocalFileStore) GetReaderAt(path string) (io.ReaderAt, int64, error) {
	path = r.filePath(path)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, fmt.Errorf(
				"file not found: %w: %w",
				err,
				e.ErrNotFound,
			)
		}
		return nil, 0, fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	return f, info.Size(), nil
}
