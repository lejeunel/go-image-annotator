package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port        int `default:"8000"`
	Path        string
	MaxPageSize int `default:"10"`
}

func NewConfig() *Config {
	var cfg Config
	err := envconfig.Process("bookstore", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(cfg.Path) == 0 {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		cfg.Path = filepath.Join(home_dir, ".cache", "bookstore", "db.sqlite")
		log.Printf("BOOKSTORE_PATH env variable not set, using default value %v\n", cfg.Path)

	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return &cfg
}
