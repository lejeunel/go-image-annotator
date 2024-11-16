package cmd

import (
	goose "github.com/pressly/goose/v3"
	a "go-image-annotator/app"
	c "go-image-annotator/config"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	cfg := c.NewConfig()
	app := a.NewApp(cfg)

	goose.SetDialect(string(goose.DialectSQLite3))
	// goose.Up(app.DB.DB, "migrations")

	app.Run(cfg.Port)
}
