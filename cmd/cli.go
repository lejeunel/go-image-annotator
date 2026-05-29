package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	clc "github.com/lejeunel/go-image-annotator/adapters/cli/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/cli/image"
)

var (
	IngestDirectoryCmd = &cobra.Command{
		Use:   "ingest-dir [dir] [collection]",
		Short: "Ingests all image located at [dir] directory into [collection]",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			collection := args[1]
			fmt.Println("ingesting directory", dir, "into collection", collection)
			im.IngestDirectory(dir, collection)
		},
	}
)

var (
	CreateCollectionCmd = &cobra.Command{
		Use:   "create-collection [name] [description]",
		Short: "Creates a new collection with [name] and [description]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := ""
			if len(args) == 2 {
				description = args[1]
			}
			name := args[0]
			clc.Create(name, description)
		},
	}
)
