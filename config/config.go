package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DbPath      string
	MaxPageSize int `default:"10"`
}

func NewConfig() *Config {
	var cfg Config
	err := envconfig.Process("GOIMANNOTATE", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(cfg.DbPath) == 0 {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		cfg.DbPath = filepath.Join(home_dir, ".cache", "go-image-annotator", "db.sqlite")
		log.Printf("GOIMANNOTATE_DBPATH env variable not set, using default value %v\n", cfg.DbPath)

	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return &cfg
}
