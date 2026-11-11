package doctor

import (
	cfg "github.com/lejeunel/go-image-annotator/config"
)

type Request struct {
	cfg.Config
}

type ExamResult struct{}

type ExamName string

const (
	InitialAdminPassword ExamName = "initial-admin-password"
	ImageStoreURI        ExamName = "image-store-uri"
	S3ImageStoreConfig   ExamName = "s3-image-store-config"
)

type Response struct {
	Name        ExamName
	Description string
	Error       error
}
