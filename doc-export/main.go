package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "docs",
	Short: "Documentation compilation",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile documentation",
	Long:  "Compile MarkDown Documentation from <input-path> into HTML files at <output-path>",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CompileDocs(args)
	},
}

func init() {

	rootCmd.AddCommand(compileCmd)
}
