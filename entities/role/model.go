package role

type Role struct {
	Id          RoleId
	Name        string
	Description string
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
