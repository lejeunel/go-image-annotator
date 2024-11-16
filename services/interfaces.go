package services

import "context"

type KeyValueStoreClient interface {
	ValidateUri(ctx context.Context, uri string) error
	Upload(ctx context.Context, uri string, data []byte, sha256 string) error
	Download(ctx context.Context, uri string) ([]byte, error)
}
