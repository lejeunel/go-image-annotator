package cmd

import (
	"context"
	tmpl "go-image-annotator/templates"
	"os"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {

	component := tmpl.Hello("Bob")
	component.Render(context.Background(), os.Stdout)

}
