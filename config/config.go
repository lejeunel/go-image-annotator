package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	SQLiteDBPath                         string   `required:"true"`
	ArtefactDir                          string   `required:"true"`
	InitialAdminEmail                    string   `required:"true"`
	InitialAdminPassword                 string   `required:"true"`
	URL                                  string   `required:"true"`
	AllowedImageFormats                  []string `default:"jpeg,png"`
	DefaultPageSize                      int      `default:"20"`
	ApiTokenLength                       int      `default:"32"`
	RandomPasswordLength                 int      `default:"10"`
	ForgotPasswordTokenExpirationMinutes int      `default:"30"`
	PasswordMinEntropy                   int      `default:"50"`
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
