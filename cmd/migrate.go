package cmd

import (
	"log/slog"
	"os"

	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Inspect and apply database migrations",
}

func init() {
	MigrateCmd.AddCommand(upCmd, downCmd, statusCmd)
}

func NewCLIDBManager() (*db.SQLiteDBManager, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	m, err := db.NewDBManagerFromEnv(logger)
	if err != nil {
		return nil, err
	}
	return m, nil
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply the next pending migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := NewCLIDBManager()
		if err != nil {
			panic(err)
		}
		manager.Up(cmd.Context())
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recently applied migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := NewCLIDBManager()
		if err != nil {
			panic(err)
		}

		manager.Down(cmd.Context())
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List migrations and whether they have been applied",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := NewCLIDBManager()
		if err != nil {
			panic(err)
		}

		manager.Status(cmd.Context())
	},
}
