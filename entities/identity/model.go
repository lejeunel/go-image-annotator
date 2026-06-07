package identity

type contextKey string

const UserKey contextKey = "user"

type Identity struct {
	Id     string
	Groups []string
	Roles  []string
}
