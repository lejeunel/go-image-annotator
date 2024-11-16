package cmd

import (
	migrationCmd "github.com/pivaldi/db-migration/cmds"
	migrationCfg "github.com/pivaldi/db-migration/config"
	"github.com/spf13/cobra"
	c "go-image-annotator/config"
)

var migrateCmd = &cobra.Command{
	Use:              "migration",
	Short:            "Database migration",
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		applyMigrationConfig()
	},
}

func applyMigrationConfig() {
	myCfg := c.NewConfig()
	dbCfg := migrationCfg.ConfigT{DBConnection: migrationCfg.DBConnectionT{DSN: myCfg.Path,
		Driver: "sqlite3"}, DBMigration: migrationCfg.DBMigrationT{Dir: "migrations"}}
	migrationCmd.SetConfig(dbCfg)

}

func init() {

	migrateCmd.AddCommand(migrationCmd.Root.Commands()...)
}
