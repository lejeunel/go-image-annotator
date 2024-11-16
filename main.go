package main

import (
	"embed"
	"github.com/pressly/goose/v3"
	"go-image-annotator/cmd"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	goose.SetBaseFS(embedMigrations)
	cmd.Execute()
}
