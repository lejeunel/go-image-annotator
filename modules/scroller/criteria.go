package scroller

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ScrollingCriteria struct {
	Collection *string
	im.Ordering
}

type Option func(*ScrollingCriteria)

func WithCollection(collection string) Option {
	return func(c *ScrollingCriteria) {
		c.Collection = &collection
	}
}

func WithOrdering(o im.Ordering) Option {
	return func(c *ScrollingCriteria) {
		c.Ordering = o
	}
}

func NewCriteria(opts ...Option) ScrollingCriteria {
	c := &ScrollingCriteria{}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}
