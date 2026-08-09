package file_store

import (
	"io"
)

type FileStore interface {
	Store(string, io.Reader) error
	Delete(string) error
	Get(string) (io.Reader, error)
}
