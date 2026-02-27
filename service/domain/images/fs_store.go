package images

import (
	"context"
	e "datahub/errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type FSKVStoreClient struct {
	RootPath string
}

func NewFSKVStoreClient(rootPath string) KeyValueStoreClient {
	return &FSKVStoreClient{RootPath: rootPath}
}

func (s *FSKVStoreClient) Scheme() string {
	return "file"
}

func (s *FSKVStoreClient) Root() string {
	return s.RootPath
}

func (s *FSKVStoreClient) ValidateUri(ctx context.Context, uri string) error {
	return nil
}

func (s *FSKVStoreClient) Upload(ctx context.Context, uri string, data []byte, sha256 string) error {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}
	path := u.Path
	parts := strings.Split(path, "/")
	dir := strings.Join(parts[:len(parts)-1], "/")

	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("creating directory: %w: %w", err, e.ErrFileStorage)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w: %w", err, e.ErrFileStorage)
	}

	return nil
}

func (s *FSKVStoreClient) Download(ctx context.Context, uri string) ([]byte, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return []byte{}, fmt.Errorf("parsing uri %v: %w: %w", uri, err, e.ErrFileStorage)
	}
	bytes, err := os.ReadFile(url.Path)
	if err != nil {
		return []byte{}, fmt.Errorf("reading file %v: %w: %w", url.Path, err, e.ErrFileStorage)
	}

	return bytes, nil
}
