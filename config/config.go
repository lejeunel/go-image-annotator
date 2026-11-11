package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBPath                               string   `required:"true" split_words:"true"`
	ImagesURI                            string   `required:"true" split_words:"true"`
	S3Endpoint                           string   `                split_words:"true"`
	S3Region                             string   `                split_words:"true"`
	S3AccessKey                          string   `                split_words:"true"`
	S3Secret                             string   `                split_words:"true"`
	S3UsePathStyle                       bool     `                split_words:"true"`
	InitialAdminEmail                    string   `required:"true" split_words:"true"`
	InitialAdminPassword                 string   `required:"true" split_words:"true"`
	URL                                  string   `                split_words:"true" default:"localhost"`
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
	SessionSecret                        string   `                split_words:"true"`
	GoogleClientId                       string   `                split_words:"true"`
	GoogleClientSecret                   string   `                split_words:"true"`
	GithubClientId                       string   `                split_words:"true"`
	GithubClientSecret                   string   `                split_words:"true"`
}

// ParseS3Config builds the S3 bucket configuration by merging the GOIA_S3_*
// variables and IMAGES_URI with the AWS SDK default chain.
//
// Credentials come from the default chain (env, shared profiles, SSO, IAM
// roles) unless both S3AccessKey and S3Secret are set, in which case they
// override it. A non-empty S3Region overrides the chain's region as well.
func (c Config) ParseS3Config(ctx context.Context) (*S3Config, error) {
	u, err := url.Parse(c.ImagesURI)
	if err != nil {
		return nil, fmt.Errorf("invalid IMAGE_URI: %w", err)
	}

	opts := []func(*awscfg.LoadOptions) error{}
	if c.S3Region != "" {
		opts = append(opts, awscfg.WithRegion(c.S3Region))
	}
	staticCredentials := c.S3AccessKey != "" && c.S3Secret != ""
	if staticCredentials {
		opts = append(opts, awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.S3AccessKey, c.S3Secret, ""),
		))
	}

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	return &S3Config{
		Bucket:            u.Host,
		Prefix:            strings.Trim(u.Path, "/"),
		Endpoint:          c.S3Endpoint,
		UsePathStyle:      c.S3UsePathStyle,
		StaticCredentials: staticCredentials,
		AWS:               awsCfg,
	}, nil
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
	Bucket            string
	Prefix            string
	Endpoint          string
	UsePathStyle      bool
	StaticCredentials bool
	AWS               aws.Config
}
