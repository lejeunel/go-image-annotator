package authorizer

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type MethodName = string

type Policy struct {
	Role        rl.RoleName
	Methods     []string
	Description string
}

type Policies []Policy

var DefaultPolicies = Policies{
	{
		"annotator",
		[]string{"Annotate", "AddMetadata", "UpdateMetadata", "DeleteMetadata"},
		"Annotate and add meta-data",
	},
	{"image-contributor", []string{
		"IngestImage",
		"ImportImage",
		"CreateCollection",
		"CloneCollection",
		"DeleteCollection",
	}, "Manage collections and ingest images"},
	{"admin", []string{"*"}, "Can do anything"},
}
