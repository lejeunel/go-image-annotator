package policy

type Policies map[string][]string

var DefaultPolicyFileName = "policies.yaml"

var DefaultPolicies = Policies{
	"viewer":            {},
	"annotator":         {"Annotate"},
	"image-contributor": {"IngestImage", "ImportImage", "CreateCollection", "CloneCollection", "DeleteCollection"},
	"admin":             {"*"},
}
