package file_store

import (
	"fmt"
	"log/slog"

	cfg "github.com/lejeunel/go-image-annotator/config"
)

func Build(cfg cfg.Config, logger slog.Logger) (FileStore, error) {
	useS3, err := cfg.DefinesS3Store()
	if err != nil {
		panic(err)
	}
	var imageStore FileStore
	if useS3 {
		s3cfg, err := cfg.S3Config()
		if err != nil {
			panic(err)
		}
		imageStore = NewS3FileStore(
			s3cfg.Endpoint,
			s3cfg.Region,
			s3cfg.AccessKey,
			s3cfg.Secret, s3cfg.Bucket, s3cfg.Prefix)
		logger.Info("using S3 image store", "bucket", s3cfg.Bucket, "prefix", s3cfg.Prefix)
	} else {
		path := fmt.Sprintf("%v/%v", cfg.LocalArtefactPath, "images")
		logger.Info("using local file system image store", "path", path)
		imageStore = NewLocalFileStore(path)
	}
	return imageStore, nil
}
