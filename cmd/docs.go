package cmd

import (
	d "datahub/docs"
	"github.com/spf13/cobra"
)

var documentationCmd = &cobra.Command{
	Use:   "docs",
	Short: "Compile documentation",
	Long:  "Parse and convert documentation.",
	Run: func(cmd *cobra.Command, args []string) {
		d.MakeDocs()
	},
}
