package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	SQLiteDBPath         string   `required:"true"`
	ArtefactDir          string   `required:"true"`
	InitialAdminEmail    string   `required:"true"`
	InitialAdminPassword string   `required:"true"`
	AllowedImageFormats  []string `default:"jpeg,png"`
	DefaultPageSize      int      `default:"20"`
	TokenLength          int      `default:"32"`
	RandomPasswordLength int      `default:"10"`
	APIPath              string   `default:"api"`
	RepoURL              string   `default:"https://github.com/lejeunel/go-image-annotator"`
	DocsURL              string   `default:"https://lejeunel.github.io/go-image-annotator/"`
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
