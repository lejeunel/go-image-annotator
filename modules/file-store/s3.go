package file_store

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileStore struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3FileStore(
	endpoint string,
	region string,
	accessKeyID string,
	secretAccessKey string,
	bucket string,
	prefix string,
) *S3FileStore {
	cfg := aws.Config{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &S3FileStore{
		client: client,
		bucket: bucket,
		prefix: prefix,
	}
}

func (s *S3FileStore) key(name string) string {
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

func (s *S3FileStore) Get(name string) (io.ReadCloser, error) {
	key := s.key(name)
	result, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("get %q: %w", name, err)
	}

	return result.Body, nil
}
