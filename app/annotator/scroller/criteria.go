package scroller

type ScrollingCriteria struct {
	Collection *string
}

type Option func(*ScrollingCriteria)

func WithCollection(collection string) Option {
	return func(c *ScrollingCriteria) {
		c.Collection = &collection
	}
}

func NewCriteria(opts ...Option) ScrollingCriteria {
	c := &ScrollingCriteria{}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}
