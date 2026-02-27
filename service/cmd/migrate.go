package main

import (
	a "datahub/app"
	m "datahub/app/migrations"
	c "datahub/config"
	clk "github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"os"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Run database migrations in the specified direction.

Valid directions:
  up    - Apply migrations
  down  - Revert migrations`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"up", "down"},
	Run: func(cmd *cobra.Command, args []string) {
		Migrate(args)
	},
}

func Migrate(args []string) {
	cfg := c.NewConfig()
	clock := clk.NewFakeClock()
	app, _, migrationProvider, ctx := a.NewApp(cfg, clock, 1)

	switch args[0] {
	case "up":
		if err := m.ApplyMigrations(ctx, migrationProvider, "up", app.Logger); err != nil {
			app.Logger.Error("Error applying up migrations", "error", err.Error())
			os.Exit(1)
		}
	case "down":
		if err := m.ApplyMigrations(ctx, migrationProvider, "down", app.Logger); err != nil {
			app.Logger.Error("Error applying down migration", "error", err.Error())
			os.Exit(1)
		}
	default:
		app.Logger.Error("Unknown migration command")
	}

}
