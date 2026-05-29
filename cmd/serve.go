package cmd

import (
	"github.com/lejeunel/go-image-annotator/site"
	"github.com/spf13/cobra"
)

var (
	port     int
	ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			site.Serve(port)
		},
	}
)

func init() {
	ServeCmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
}
