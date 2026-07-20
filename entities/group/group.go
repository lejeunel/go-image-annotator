package group

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type Group struct {
	Id          GroupId
	Name        string
	Description string
}

func NewGroup(id GroupId, name string, opts ...Option) Group {
	c := &Group{Id: id, Name: name}
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

type GroupId struct {
	uuidw.UUIDWrapper[GroupId]
}

func NewGroupId() GroupId {
	return GroupId{uuidw.UUIDWrapper[GroupId]{UUID: uuid.New()}}
}
