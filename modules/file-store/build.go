package file_store

import (
	"log/slog"

	cfg "github.com/lejeunel/go-image-annotator/config"
)

func Build(cfg cfg.Config, logger slog.Logger) (FileStore, error) {
	useS3, err := cfg.DefinesS3Store()
	if err != nil {
		return nil, err
	}
	if !useS3 {
		path, err := cfg.LocalImagePath()
		if err != nil {
			return nil, err
		}
		logger.Info("using local file system image store", "path", *path)
		return NewLocalFileStore(*path), nil
	}

	s3cfg, err := cfg.S3Config()
	if err != nil {
		return nil, err
	}
	imageStore, err := NewS3FileStore(
		s3cfg.Endpoint,
		s3cfg.Region,
		s3cfg.AccessKey,
		s3cfg.Secret, s3cfg.Bucket, s3cfg.Prefix, s3cfg.UsePathStyle)
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
		s3cfg.Region,
		"path_style",
		s3cfg.UsePathStyle,
		"credentials",
		credentialSource(s3cfg.AccessKey, s3cfg.Secret),
	)
	return imageStore, nil
}

func credentialSource(accessKey, secret string) string {
	if accessKey != "" && secret != "" {
		return "static"
	}
	return "default AWS chain"
}
