package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	SQLiteDBPath        string   `required:"true"`
	ArtefactDir         string   `required:"true"`
	AllowedImageFormats []string `default:"jpeg,png"`
	DefaultPageSize     int      `default:"10"`
	TokenLength         int      `default:"32"`
	APIPath             string   `default:"api"`
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
