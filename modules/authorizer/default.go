package authorizer

var DefaultPolicies = Policies{
	"viewer":            {},
	"annotator":         {"Annotate"},
	"image-contributor": {"IngestImage", "ImportImage", "CreateCollection", "CloneCollection", "DeleteCollection"},
	"admin":             {"*"},
}
