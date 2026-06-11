package image

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	IngestDirCmd = &cobra.Command{
		Use:   "ingest-dir [dir] [collection]",
		Short: "Ingests all image located at [dir] directory into [collection]",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			collection := args[1]
			fmt.Println("ingesting directory", dir, "into collection", collection)
			IngestDirectory(dir, collection)
		},
	}
)
