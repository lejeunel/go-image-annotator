package collection

import (
	"github.com/google/uuid"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
	"time"
)

type Collection struct {
	Id          CollectionId
	Name        string
	Description string
	CreatedAt   time.Time
	Group       *g.Group
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

func WithGroup(g g.Group) Option {
	return func(c *Collection) {
		c.Group = &g
	}
}

type UpdateModel struct {
	Name           string
	NewName        string
	NewDescription string
}

type CollectionId struct {
	uuidw.UUIDWrapper[CollectionId]
}

func NewCollectionId() CollectionId {
	return CollectionId{uuidw.UUIDWrapper[CollectionId]{UUID: uuid.New()}}
}
