package cmd

import (
	"github.com/spf13/cobra"

	clc "github.com/lejeunel/go-image-annotator/adapters/cli/collection"
)

var (
	group               string
	description         string
	CreateCollectionCmd = &cobra.Command{
		Use:   "create-collection [name] [group] [description]",
		Short: "Creates a new collection with [name] in [group] with [description]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			var groupPtr *string
			if cmd.Flags().Changed("group") {
				groupPtr = &group
			}
			clc.Create(name, groupPtr, description)
		},
	}
)

func init() {
	CreateCollectionCmd.Flags().StringVarP(&group, "group", "g", "", "an optional group")
	CreateCollectionCmd.Flags().StringVarP(&description, "description", "d", "", "an optional description")
}
