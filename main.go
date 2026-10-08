package main

import (
	"fmt"
	"os"

	"github.com/lejeunel/go-image-annotator/migrate"
	"github.com/lejeunel/go-image-annotator/server"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-image-annotator",
	Short: "Image annotation platform",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(migrate.Cmd)
}
