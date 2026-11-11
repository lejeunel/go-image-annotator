package file_store

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/lejeunel/go-image-annotator/config"
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

	if u.Scheme == "s3" {
		return S3FileStoreType, u.Host, nil
	}
	return LocalFSFileStoreType, u.Path, nil
}

func Build(cfg config.Config, logger slog.Logger) (FileStore, error) {
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
		logger.Info(
			"using S3 image store",
			"bucket",
			s3cfg.Bucket,
			"prefix",
			s3cfg.Prefix,
			"endpoint",
			s3cfg.Endpoint,
			"region",
			s3cfg.AWS.Region,
			"path_style",
			s3cfg.UsePathStyle,
			"credentials",
			credentialSource(s3cfg.StaticCredentials),
		)
		return NewS3FileStore(*s3cfg), nil

	case LocalFSFileStoreType:
		logger.Info("using local file system image store", "path", path)
		return NewLocalFileStore(path), nil
	}
	return nil, fmt.Errorf("unsupported image store type: %v", storeType)
}

func credentialSource(static bool) string {
	if static {
		return "static"
	}
	return "default AWS chain"
}
