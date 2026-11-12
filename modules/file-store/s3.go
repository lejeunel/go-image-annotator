package file_store

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	cfg "github.com/lejeunel/go-image-annotator/config"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type S3FileStore struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3FileStore builds a store backed by an S3-compatible service.
func NewS3FileStore(c cfg.S3Config) *S3FileStore {
	client := s3.NewFromConfig(c.AWS, func(o *s3.Options) {
		if c.Endpoint != "" {
			o.BaseEndpoint = aws.String(c.Endpoint)
		}
		o.UsePathStyle = c.UsePathStyle
	})

	return &S3FileStore{
		client: client,
		bucket: c.Bucket,
		prefix: c.Prefix,
	}
}

func (s *S3FileStore) key(name string) string {
	if s.prefix == "" {
		return name
	}
	return s.prefix + "/" + name
}

func (s *S3FileStore) Store(name string, r io.Reader) error {
	key := s.key(name)
	_, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("store %q: %w", name, err)
	}

	return nil
}

func (s *S3FileStore) Delete(name string) error {
	key := s.key(name)
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("delete %q: %w", name, err)
	}

	return nil
}

// Get returns a byte reader given an input path.
// Returns an error when the requested file does not exist.
func (s *S3FileStore) Get(name string) (io.ReadCloser, error) {
	key := s.key(name)
	result, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, fmt.Errorf("get %q: %w: %w", name, err, e.ErrNotFound)
		}
		return nil, fmt.Errorf("get %q: %w: %w", name, err, e.ErrInternal)
	}

	return result.Body, nil
}
