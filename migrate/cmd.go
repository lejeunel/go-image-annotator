package migrate

import (
	"log/slog"
	"os"

	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Inspect and apply database migrations",
}

func init() {
	Cmd.AddCommand(upCmd, downCmd, statusCmd)
}

func NewCLIDBManager() *db.SQLiteDBManager {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	m := db.NewDBManagerFromEnv(logger)
	return &m
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply the next pending migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		NewCLIDBManager().Up(cmd.Context())
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recently applied migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		NewCLIDBManager().Down(cmd.Context())
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List migrations and whether they have been applied",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		NewCLIDBManager().Status(cmd.Context())
	},
}
