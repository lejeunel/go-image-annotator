package auth

type AuthRule struct {
	Method      string
	IgnoreGroup bool
	Roles       []string
	AdminOnly   bool
}
