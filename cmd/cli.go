package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	clc "github.com/lejeunel/go-image-annotator/adapters/cli/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/cli/image"
	usr "github.com/lejeunel/go-image-annotator/adapters/cli/user"
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
		Use:   "create-collection [name] [group] [description]",
		Short: "Creates a new collection with [name] in [group] with [description]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			description := ""
			if len(args) == 3 {
				description = args[2]
			}
			name := args[0]
			group := args[1]
			clc.Create(name, &group, description)
		},
	}
)

var (
	isAdmin       bool
	CreateUserCmd = &cobra.Command{
		Use:   "create-user [id] [is-admin]",
		Short: "Creates a new user with [id] and [group] with [description]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			usr.Create(args[0], isAdmin)
		},
	}
)

func init() {
	CreateUserCmd.Flags().BoolVar(&isAdmin, "is-admin", false, "set admin privileges")
}
