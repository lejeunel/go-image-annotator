package file_store

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lejeunel/go-image-annotator/config"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type FileStoreType int

const (
	UnknownFileStoreType FileStoreType = iota
	LocalFSFileStoreType
	S3FileStoreType
)

var fileStoreTypeName = map[FileStoreType]string{
	UnknownFileStoreType: "unknown",
	LocalFSFileStoreType: "local file system",
	S3FileStoreType:      "s3 bucket",
}

func (t FileStoreType) String() string {
	return fileStoreTypeName[t]
}

// ParseURI parses a URI and decides which file-store to use.
// Also returns the store location as a parsed string
func ParseURI(uri string) (FileStoreType, string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return UnknownFileStoreType, "", fmt.Errorf("invalid IMAGE_URI: %w", err)
	}

	switch u.Scheme {
	case "s3":
		return S3FileStoreType, u.Host, nil
	case "file":
		return LocalFSFileStoreType, u.Path, nil
	}
	return UnknownFileStoreType, "", fmt.Errorf(
		"validating prefix %v in URI %v: %w",
		u.Scheme,
		uri,
		e.ErrValidation,
	)
}

func Build(cfg config.Config) (FileStore, error) {
	storeType, path, err := ParseURI(cfg.ImagesURI)
	if err != nil {
		return nil, err
	}
	switch storeType {
	case S3FileStoreType:
		s3cfg, err := cfg.ParseS3Config(context.Background())
		if err != nil {
			return nil, err
		}
		return NewS3FileStore(*s3cfg), nil

	case LocalFSFileStoreType:
		return NewLocalFileStore(path)
	}
	return nil, fmt.Errorf("unsupported image store type: %v", storeType)
}
