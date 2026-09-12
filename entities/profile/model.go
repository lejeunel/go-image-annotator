package role

type Profile struct {
	Id          ProfileId
	Name        string
	Description string
	Labels      []string
}

func NewProfile(id ProfileId, name string, opts ...Option) Profile {
	r := &Profile{Id: id, Name: name}
	for _, opt := range opts {
		opt(r)
	}
	return *r
}

type Option func(*Profile)

func WithDescription(d string) Option {
	return func(r *Profile) {
		r.Description = d
	}
}

func WithLabels(labels []string) Option {
	return func(r *Profile) {
		r.Labels = labels
	}
}

type UpdatableModel struct {
	Name           string
	NewName        string
	NewDescription string
	NewLabels      []string
}
