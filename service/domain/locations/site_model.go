package locations

import (
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type SiteId struct{ g.UUIDWrapper[SiteId] }

type SiteOption func(*Site)

func WithGroupOption(group string) SiteOption {
	return func(s *Site) {
		s.Group = group
	}
}

func NewSiteId() *SiteId {
	id := uuid.New()
	return &SiteId{g.UUIDWrapper[SiteId]{UUID: id}}

}

func (id *SiteId) Equal(to *SiteId) bool {
	if (to == nil) && (id == nil) {
		return true
	}
	if id.String() == to.String() {
		return true
	}
	return false
}

func NewSiteIdFromUUID(id uuid.UUID) *SiteId {
	return &SiteId{g.UUIDWrapper[SiteId]{UUID: id}}
}

func NewSiteIdFromString(s string) (*SiteId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing site id: %w", e.ErrValidation)
	}
	return NewSiteIdFromUUID(id), nil

}

type Site struct {
	Id        SiteId
	Name      string
	Group     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSite(name string, opts ...SiteOption) (*Site, error) {
	site := &Site{Id: *NewSiteId(), Name: name}

	for _, opt := range opts {
		opt(site)
	}
	if err := site.Validate(); err != nil {
		return nil, err
	}
	return site, nil
}

func (s Site) Validate() error {
	if err := validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Name, validation.Match(g.ResourceNameRegExp)),
	); err != nil {
		return fmt.Errorf("validating site name (%v): %w", s.Name, e.ErrResourceName)
	}
	return nil
}
