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
	AllowedImageMIMETypes                []string `default:"image/jpeg,image/png" split_words:"true"`
	DefaultPageSize                      int      `default:"20" split_words:"true"`
	ApiTokenLength                       int      `default:"32" split_words:"true"`
	RandomPasswordLength                 int      `default:"10" split_words:"true"`
	ForgotPasswordTokenExpirationMinutes int      `default:"30" split_words:"true"`
	PasswordMinEntropy                   int      `default:"50" split_words:"true"`
	MaxNumTasksPerUser                   int      `default:"50" split_words:"true"`
	MaxArchiveMB                         int      `default:"500" split_words:"true"`
	SMTPUsername                         string   `split_words:"true"`
	SMTPPassword                         string   `split_words:"true"`
	SMTPHost                             string   `split_words:"true"`
	SMTPPort                             int      `split_words:"true"`
	GoogleClientId                       string   `split_words:"true"`
	GoogleClientSecret                   string   `split_words:"true"`
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
