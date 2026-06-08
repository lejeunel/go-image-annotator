package group

type Group struct {
	Name        string
	Description string
}

func NewGroup(name string, opts ...Option) Group {
	c := &Group{Name: name}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

type Option func(*Group)

func WithDescription(d string) Option {
	return func(c *Group) {
		c.Description = d
	}
}
