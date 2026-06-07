package user

type User struct {
	Id      string
	HashPAT []byte
	Roles   []string
	Groups  []string
}

func NewUser(id string, opts ...Option) User {
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
