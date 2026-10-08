package migrate

import (
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Inspect and apply database migrations",
	Long: "Inspect and apply database migrations.\n\n" +
		"The server applies every pending migration on startup, so these commands are\n" +
		"meant for local testing: checking what is pending, or stepping a single\n" +
		"migration up and down while writing it.",
}

func init() {
	Cmd.AddCommand(upCmd, downCmd, statusCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply the next pending migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sqlite.NewDBManagerFromEnv().Up(cmd.Context())
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recently applied migration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sqlite.NewDBManagerFromEnv().Down(cmd.Context())
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List migrations and whether they have been applied",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sqlite.NewDBManagerFromEnv().Status(cmd.Context())
	},
}
