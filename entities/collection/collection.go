package collection

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
	"time"
)

type CollectionName = string

type Collection struct {
	Id          CollectionId
	Name        string
	Description string
	CreatedAt   time.Time
	Group       *string
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
		c.Group = &g
	}
}

type UpdateModel struct {
	Name           string
	NewName        string
	NewDescription string
	NewGroup       *string
}

type CollectionId struct {
	uuidw.UUIDWrapper[CollectionId]
}

func NewCollectionId() CollectionId {
	return CollectionId{uuidw.UUIDWrapper[CollectionId]{UUID: uuid.New()}}
}
