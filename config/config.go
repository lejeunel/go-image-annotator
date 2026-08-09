package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ArtefactPath                         string   `required:"true"`
	InitialAdminEmail                    string   `required:"true"`
	InitialAdminPassword                 string   `required:"true"`
	URL                                  string   `required:"true"`
	AllowedImageMIMETypes                []string `default:"image/jpeg,image/png"`
	DefaultPageSize                      int      `default:"20"`
	ApiTokenLength                       int      `default:"32"`
	RandomPasswordLength                 int      `default:"10"`
	ForgotPasswordTokenExpirationMinutes int      `default:"30"`
	PasswordMinEntropy                   int      `default:"50"`
	MaxNumTasksPerUser                   int      `default:"50"`
	SMTPUsername                         string
	SMTPPassword                         string
	SMTPHost                             string
	SMTPPort                             int
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
