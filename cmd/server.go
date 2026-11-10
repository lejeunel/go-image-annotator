package cmd

import (
	"github.com/lejeunel/go-image-annotator/server"
	"github.com/spf13/cobra"
)

var (
	port      int
	ServerCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Serve(port)
		},
	}
)

func init() {
	ServerCmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
}
