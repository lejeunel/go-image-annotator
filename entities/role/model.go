package role

type RoleName = string

type Role struct {
	Id          RoleId
	Name        RoleName
	Description string
	Methods     []string
}

func NewRole(id RoleId, name string, opts ...Option) Role {
	r := &Role{Id: id, Name: name}
	for _, opt := range opts {
		opt(r)
	}
	return *r
}

type Option func(*Role)

func WithDescription(d string) Option {
	return func(r *Role) {
		r.Description = d
	}
}

func WithMethods(methods []string) Option {
	return func(r *Role) {
		r.Methods = methods
	}
}

type UpdatableModel struct {
	Name           string
	NewName        string
	NewDescription string
	NewMethods     []string
}
