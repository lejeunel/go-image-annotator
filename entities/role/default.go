package role

type DefaultRole struct {
	Name        string
	Description string
}

var DefaultRoleNames = []DefaultRole{
	{"annotator", "can annotate images"},
	{"image-contributor", "can create collections and add images"},
	{"admin", "can do anything"},
}
