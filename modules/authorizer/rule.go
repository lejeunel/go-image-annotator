package authorizer

type AuthRule struct {
	Method      string
	IgnoreGroup bool
	Roles       []string
	AdminOnly   bool
}
