package images

import (
	"context"
)

type MockKVStoreClient struct {
	items map[string][]byte
}

func NewMockKVStoreClient() KeyValueStoreClient {
	return &MockKVStoreClient{items: make(map[string][]byte)}
}

func (s *MockKVStoreClient) Scheme() string {
	return "scheme"
}

func (s *MockKVStoreClient) Root() string {
	return "mybucket"
}

func (s *MockKVStoreClient) ValidateUri(ctx context.Context, uri string) error {
	return nil
}

func (s *MockKVStoreClient) Upload(ctx context.Context, uri string, data []byte, sha256 string) error {
	s.items[uri] = data
	return nil
}

func (s *MockKVStoreClient) Download(ctx context.Context, uri string) ([]byte, error) {

	data := s.items[uri]
	return data, nil
}
