package main

import (
	"fmt"
	mdxp "gomdxp/internal"
	"os"
)

func Compile(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("requires exactly 2 arguments: source and destination")
	}

	if _, err := os.Stat(args[0]); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}

	builder := mdxp.NewBuilder(mdxp.NewGoldMarkParser())
	if err := mdxp.ParseMarkDownToHTML(builder, os.DirFS(args[0])); err != nil {
		return fmt.Errorf("parsing markdown pages: %w", err)
	}
	fmt.Println("Successfully parsed MarkDown pages.")

	if err := mdxp.ExportPagesToHTML(builder, args[1]); err != nil {
		return fmt.Errorf("exporting pages to HTML: %w", err)
	}
	fmt.Printf("Successfully exported pages to HTML at %v\n", args[1])

	return nil

}
