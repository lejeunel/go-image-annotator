package collections

import (
	pro "datahub/domain/annotation_profiles"
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type CollectionId struct{ g.UUIDWrapper[CollectionId] }

func NewCollectionIdFromUUID(id uuid.UUID) *CollectionId {
	return &CollectionId{g.UUIDWrapper[CollectionId]{UUID: id}}
}

func NewCollectionIdFromString(s string) (*CollectionId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing collection id: %v: %w", s, e.ErrValidation)
	}
	return NewCollectionIdFromUUID(id), nil

}

func NewCollectionId() *CollectionId {
	id := uuid.New()
	return &CollectionId{g.UUIDWrapper[CollectionId]{UUID: id}}

}

type Collection struct {
	Id          CollectionId
	Name        string
	Description string
	Group       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProfileId   *pro.AnnotationProfileId

	Profile *pro.AnnotationProfile
}

type CollectionOption func(*Collection)

func WithDescription(desc string) CollectionOption {
	return func(c *Collection) {
		c.Description = desc
	}
}

func WithGroup(group string) CollectionOption {
	return func(c *Collection) {
		c.Group = group
	}
}

func New(name string, opts ...CollectionOption) (*Collection, error) {
	id := NewCollectionId()
	collection := &Collection{Id: *id, Name: name}
	for _, opt := range opts {
		opt(collection)
	}
	if err := collection.Validate(); err != nil {
		return nil, err
	}
	return collection, nil
}

func (s Collection) Validate() error {
	if err := validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Name, validation.Match(g.ResourceNameRegExp)),
	); err != nil {
		return fmt.Errorf("validating collection name (%v): %w", s.Name, e.ErrResourceName)
	}
	return nil
}
