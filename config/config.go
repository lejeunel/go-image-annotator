package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ArtefactPath                         string   `required:"true" split_words:"true"`
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

func Parse() Config {
	var cfg Config
	err := envconfig.Process("GOIA", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}

type APIConfig struct {
	APIPath          string
	APIDocsPath      string
	OpenAPISpecsPath string
}
