package role

const AdminRoleName = "admin"

type DefaultRole struct {
	Name        string
	Description string
}

var DefaultRoleNames = []DefaultRole{
	{"annotator", "can annotate images"},
	{"image-contributor", "can create collections and add images"},
	{AdminRoleName, "can do anything"},
}
