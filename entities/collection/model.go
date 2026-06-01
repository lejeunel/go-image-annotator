package collection

import "time"

type Collection struct {
	Id          CollectionId
	Name        string
	Description string
	CreatedAt   time.Time
	Group       string
}

func NewCollection(id CollectionId, name string, opts ...Option) Collection {
	c := &Collection{Id: id, Name: name}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

type Option func(*Collection)

func WithDescription(d string) Option {
	return func(c *Collection) {
		c.Description = d
	}
}

func WithCreatedAt(t time.Time) Option {
	return func(c *Collection) {
		c.CreatedAt = t
	}
}

func WithGroup(g string) Option {
	return func(c *Collection) {
		c.Group = g
	}
}
