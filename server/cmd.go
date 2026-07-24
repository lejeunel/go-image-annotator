package server

import (
	"github.com/spf13/cobra"
)

var (
	port          int
	authRulesPath string
	Cmd           = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			handler := Make("http://localhost", port)
			Serve(handler, port)
		},
	}
)

func init() {
	Cmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
}
