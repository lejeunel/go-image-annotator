package cmd

import (
	site "datahub/site"
	"github.com/spf13/cobra"
)

var (
	port     int
	migrate  bool
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			site.Serve(port, migrate)
		},
	}
)

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
	serveCmd.Flags().BoolVarP(&migrate, "migrate", "m", false, "apply DB migrations")
}
