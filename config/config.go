package config

import (
	"fmt"
	"log"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LocalArtefactPath                    string   `required:"true" split_words:"true"`
	ImagesURI                            string   `required:"true" split_words:"true"`
	S3Endpoint                           string   `                split_words:"true"`
	S3Prefix                             string   `                split_words:"true"`
	S3Region                             string   `                split_words:"true"`
	S3AccessKey                          string   `                split_words:"true"`
	S3Secret                             string   `                split_words:"true"`
	InitialAdminEmail                    string   `required:"true" split_words:"true"`
	InitialAdminPassword                 string   `required:"true" split_words:"true"`
	URL                                  string   `required:"true" split_words:"true"`
	AllowedImageMIMETypes                []string `                split_words:"true" default:"image/jpeg,image/png"`
	DefaultPageSize                      int      `                split_words:"true" default:"20"`
	MaxPageSize                          int      `                split_words:"true" default:"50"`
	ApiTokenLength                       int      `                split_words:"true" default:"32"`
	RandomPasswordLength                 int      `                split_words:"true" default:"10"`
	ForgotPasswordTokenExpirationMinutes int      `                split_words:"true" default:"30"`
	PasswordMinEntropy                   int      `                split_words:"true" default:"50"`
	MaxNumTasksPerUser                   int      `                split_words:"true" default:"50"`
	MaxArchiveMB                         int      `                split_words:"true" default:"500"`
	SMTPUsername                         string   `                split_words:"true"`
	SMTPPassword                         string   `                split_words:"true"`
	SMTPHost                             string   `                split_words:"true"`
	SMTPPort                             int      `                split_words:"true"`
	GoogleClientId                       string   `                split_words:"true"`
	GoogleClientSecret                   string   `                split_words:"true"`
}

func (c Config) S3Config() (*S3Config, error) {
	u, err := url.Parse(c.ImagesURI)
	if err != nil {
		return nil, fmt.Errorf("invalid IMAGE_URI: %w", err)
	}

	return &S3Config{
		Bucket:    u.Host,
		Prefix:    u.Path,
		Endpoint:  c.S3Endpoint,
		Region:    c.S3Region,
		AccessKey: c.S3AccessKey,
		Secret:    c.S3AccessKey,
	}, nil
}

func (c Config) DefinesS3Store() (bool, error) {
	u, err := url.Parse(c.ImagesURI)
	if err != nil {
		return false, fmt.Errorf("invalid IMAGE_URI: %w", err)
	}
	if u.Scheme == "s3" {
		return true, nil
	}
	return false, nil
}

// Parse builds the main Config struct
// from environment variables
func Parse() Config {
	var cfg Config
	err := envconfig.Process("GOIA", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}

type S3Config struct {
	Bucket    string
	Prefix    string
	Endpoint  string
	Region    string
	AccessKey string
	Secret    string
}
