package server

import (
	"github.com/lejeunel/go-image-annotator/modules/auth"
	"github.com/spf13/cobra"
)

var (
	port          int
	authSpecsPath string
	Cmd           = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			specs, err := auth.ReadAuthSpecsFromPath(authSpecsPath)
			if err != nil {
				panic(err)
			}

			authorizer, err := auth.New(*specs)
			if err != nil {
				panic(err)
			}
			handler := Make(*authorizer)
			Serve(handler, port)
		},
	}
)

func init() {
	Cmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
	Cmd.Flags().StringVarP(&authSpecsPath, "auth-specs", "a", "", "path to yaml authentication specification file")
}
