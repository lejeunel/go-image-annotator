package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3ConfigParsesBucketAndPrefix(t *testing.T) {
	for _, tc := range []struct {
		uri    string
		bucket string
		prefix string
	}{
		{"s3://my-bucket/images", "my-bucket", "images"},
		{"s3://my-bucket/nested/images", "my-bucket", "nested/images"},
		{"s3://my-bucket/images/", "my-bucket", "images"},
		{"s3://my-bucket", "my-bucket", ""},
		{"s3://my-bucket/", "my-bucket", ""},
	} {
		t.Run(tc.uri, func(t *testing.T) {
			s3cfg, err := Config{ImagesURI: tc.uri}.S3Config()
			require.NoError(t, err)
			assert.Equal(t, tc.bucket, s3cfg.Bucket)
			assert.Equal(t, tc.prefix, s3cfg.Prefix)
		})
	}
}
