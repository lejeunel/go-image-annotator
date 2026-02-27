package main

import (
	docs "docexport/internal"
	"fmt"
	"os"
)

func CompileDocs(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("requires exactly 2 arguments: source and destination")
	}

	if _, err := os.Stat(args[0]); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}

	builder := docs.NewBuilder(docs.NewGoldMarkParser())
	if err := docs.ParseMarkDownToHTML(builder, os.DirFS(args[0])); err != nil {
		return fmt.Errorf("parsing markdown pages: %w", err)
	}
	fmt.Println("Successfully parsed MarkDown pages.")

	if err := docs.ExportPagesToHTML(builder, args[1]); err != nil {
		return fmt.Errorf("exporting pages to HTML: %w", err)
	}

	fmt.Printf("Successfully exported pages to HTML at %v\n", args[1])

	return nil

}
