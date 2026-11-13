package file_store

import (
	"io"
)

type FileStore interface {
	Store(string, io.Reader) error
	Delete(string) error
	Get(string) (io.ReadCloser, error)
	// Exists checks for the file without fetching its content.
	Exists(string) (bool, error)
}
