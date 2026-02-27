package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"slices"
)

type Config struct {
	Email              string   `default:"anonymous@mail.com"`
	MaxPageSize        int      `default:"10"`
	LocalPath          string   `default:"."`
	PagerSize          int      `default:"10"`
	Mode               string   `default:"test"`
	AllowedImageTypes  []string `default:"thermal,rgb,gray"`
	PostGreSqlUser     string   `default:""`
	PostGreSqlPassword string   `default:""`
	PostGreSqlHost     string   `default:""`
	PostGreSqlPort     int      `default:""`

	TestingEntitlements []string `default:"admin"`
	SignOutURL          string   `default:"#"`

	TargetImageWidth int `default:"700"`
}

func (c *Config) DBDriver() string {
	switch c.Mode {
	case "prod":
		return "postgres"
	case "test", "dev":
		return "sqlite3"
	}
	return ""
}

func (c *Config) DBDataSourceName() string {
	switch c.Mode {
	case "prod":
		return fmt.Sprintf("postgres://%v:%v@%v:%v/datahub?sslmode=disable",
			c.PostGreSqlUser, c.PostGreSqlPassword,
			c.PostGreSqlHost, c.PostGreSqlPort)
	case "test", "dev":
		return c.LocalPath + "/db.sqlite"
	}
	return ""
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
