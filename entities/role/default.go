package role

type DefaultRole struct {
	Name        string
	Description string
}

var DefaultRoleNames = []DefaultRole{
	{"viewer", "can view images"},
	{"annotator", "can annotate images"},
	{"image-contributor", "can create collections and add images"},
	{"admin", "can do anything"},
}
