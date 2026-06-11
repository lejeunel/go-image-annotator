package collection

import (
	"github.com/spf13/cobra"
)

var (
	group       string
	description string
	CreateCmd   = &cobra.Command{
		Use:   "create-collection [name] [group] [description]",
		Short: "Creates a new collection with [name] in [group] with [description]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			var groupPtr *string
			if cmd.Flags().Changed("group") {
				groupPtr = &group
			}
			Create(name, groupPtr, description)
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&group, "group", "g", "", "an optional group")
	CreateCmd.Flags().StringVarP(&description, "description", "d", "", "an optional description")
}
