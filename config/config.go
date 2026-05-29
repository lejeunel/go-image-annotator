package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	DBPath              string   `required:"true"`
	ArtefactDir         string   `required:"true"`
	AllowedImageFormats []string `default:"jpeg,png"`
	DefaultPageSize     int      `default:"10"`
}

func Parse() Config {
	var cfg Config
	err := envconfig.Process("GOIA", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}
