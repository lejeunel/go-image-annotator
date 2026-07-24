package file_store

import (
	"io"
)

type Interface interface {
	Store(string, io.Reader) error
	Delete(string) error
	Get(string) (io.Reader, error)
}

type ReadInterface interface {
	Get(string) (io.Reader, error)
}
