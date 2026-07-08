package server

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/spf13/cobra"
)

func MakeAuthorizer(path string) auth.Authorizer {
	authorizer := auth.NewDefault()
	if authRulesPath != "" {
		rules, err := auth.ReadAuthRulesFromPath(authRulesPath)
		if err != nil {
			panic(err)
		}
		authorizer.SetAuthRules(*rules)
	}
	return authorizer

}

var (
	port          int
	authRulesPath string
	Cmd           = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			authorizer := MakeAuthorizer(authRulesPath)
			handler := Make(authorizer, "http://localhost", port)
			Serve(handler, port)
		},
	}
)

func init() {
	Cmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
	Cmd.Flags().StringVarP(&authRulesPath, "auth-rules", "a", "", "path to yaml file that specifies authorization rules")
}
