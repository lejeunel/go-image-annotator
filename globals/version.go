package globals

var (
	Version     = "dev"
	Commit      = "unknown"
	Date        = "unknown"
	RepoURL     = "https://github.com/lejeunel/go-image-annotator"
	DocsURL     = "https://lejeunel.github.io/go-image-annotator/"
	PackageName = "go-image-annotator"
)

type Info struct {
	Version string
	Commit  string
	Date    string
}
