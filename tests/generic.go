package tests

import (
	"context"
	_ "embed"
	"fmt"
	goose "github.com/pressly/goose/v3"
	a "go-image-annotator/app"
	r "go-image-annotator/repositories/sql"
	s "go-image-annotator/services"
	"testing"
)

//go:embed test-data/sample-image.png
var testImage []byte

type Services struct {
	Images *s.ImageService
}

type MockKVStoreClient struct {
	items         map[string][]byte
	AllowedScheme string
	AllowedPrefix string
}

func NewMockKVStoreClient() *MockKVStoreClient {
	return &MockKVStoreClient{items: make(map[string][]byte),
		AllowedScheme: "scheme", AllowedPrefix: "mybucket"}
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

func NewTestComponents(t *testing.T) (Services, context.Context) {
	db := a.NewSQLiteConnection(":memory:")
	goose.SetLogger(goose.NopLogger())
	goose.SetDialect(string(goose.DialectSQLite3))
	err := goose.Up(db.DB, "../migrations")
	if err != nil {
		panic(err)
	}
	imageRepo := r.NewSQLImageRepo(db)
	KVStore := NewMockKVStoreClient()

	imageService := s.ImageService{KeyValueStoreClient: KVStore, ImageRepo: imageRepo, MaxPageSize: 2,
		DefaultPageSize: 2, RemoteScheme: "scheme", RemoteBucketName: "mybucket"}

	return Services{Images: &imageService}, context.Background()

}

func AssertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(fmt.Printf("did not want an error but got: %v", err))
	}
}
