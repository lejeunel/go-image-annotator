package user

type User struct {
	Id      string
	HashPAT []byte
	Roles   []string
	Groups  []string
	IsAdmin bool
}

var UserContextKey = "user"

type UserId = string

func NewUser(id UserId, opts ...Option) User {
	l := &User{Id: id}
	for _, opt := range opts {
		opt(l)
	}
	return *l
}

type Option func(*User)

func WithHashedPersonalAccessToken(h []byte) Option {
	return func(l *User) {
		l.HashPAT = h
	}
}

func WithGroups(groups []string) Option {
	return func(l *User) {
		l.Groups = groups
	}
}

func WithAdmin(admin bool) Option {
	return func(l *User) {
		if admin {
			l.IsAdmin = true
		}
	}
}

func WithRoles(roles []string) Option {
	return func(l *User) {
		l.Roles = roles
	}
}
