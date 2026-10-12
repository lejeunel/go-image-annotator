package file_store

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileStore struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3FileStore builds a store backed by an S3-compatible service.
//
// Credentials are resolved through the AWS default chain, which covers IAM
// roles (EC2/ECS metadata, EKS web identity), SSO and shared profiles. Passing
// a non-empty accessKeyID and secretAccessKey overrides that chain, which is
// what non-AWS deployments (MinIO and friends) are expected to do.
//
// usePathStyle addresses buckets as endpoint/bucket/key instead of the
// virtual-host form bucket.endpoint/key. Services reached through a custom
// endpoint without wildcard DNS, MinIO most notably, need it.
func NewS3FileStore(
	endpoint string,
	region string,
	accessKeyID string,
	secretAccessKey string,
	bucket string,
	prefix string,
	usePathStyle bool,
) (*S3FileStore, error) {
	opts := []func(*awscfg.LoadOptions) error{}
	if region != "" {
		opts = append(opts, awscfg.WithRegion(region))
	}
	if accessKeyID != "" && secretAccessKey != "" {
		opts = append(opts, awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		))
	}

	cfg, err := awscfg.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = usePathStyle
	})

	return &S3FileStore{
		client: client,
		bucket: bucket,
		prefix: prefix,
	}, nil
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
