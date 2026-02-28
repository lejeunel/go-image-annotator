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
	Use:   "gomdxp",
	Short: "Markdown Exporter",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

var compileCmd = &cobra.Command{
	Use:   "export",
	Short: "Export pages",
	Long:  "Export MarkDown pages from <input-path> to <output-path>",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Compile(args)
	},
}

func init() {

	rootCmd.AddCommand(compileCmd)
}
