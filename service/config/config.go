package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"slices"
)

type Config struct {
	Email             string   `default:"anonymous@mail.com"`
	MaxPageSize       int      `default:"10"`
	LocalPath         string   `default:"."`
	PagerSize         int      `default:"10"`
	Mode              string   `default:"test"`
	AllowedImageTypes []string `default:"thermal,rgb,gray"`

	TestingEntitlements []string `default:"admin"`
	SignOutURL          string   `default:"#"`

	TargetImageWidth int `default:"700"`
}

func NewConfig() *Config {
	var cfg Config
	err := envconfig.Process("DATAHUB", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	allowedModes := []string{"test", "dev", "prod"}
	if !slices.Contains(allowedModes, cfg.Mode) {
		log.Fatalf("Provided mode %v not allowed. It should be one of %v",
			cfg.Mode, allowedModes)
	}

	return &cfg
}
