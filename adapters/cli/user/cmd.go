package user

import (
	"github.com/spf13/cobra"
)

var (
	isAdmin   bool
	CreateCmd = &cobra.Command{
		Use:   "create-user [id] [is-admin]",
		Short: "Creates a new user with [id] and [group] with [description]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Create(args[0], isAdmin)
		},
	}
)

func init() {
	CreateCmd.Flags().BoolVar(&isAdmin, "is-admin", false, "set admin privileges")
}
