package labels

import (
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type LabelId struct{ g.UUIDWrapper[LabelId] }

func NewLabelIdFromUUID(id uuid.UUID) *LabelId {
	return &LabelId{g.UUIDWrapper[LabelId]{UUID: id}}
}

func NewLabelIdFromString(s string) (*LabelId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing label id: %w", e.ErrValidation)
	}
	return NewLabelIdFromUUID(id), nil

}

func NewLabelId() *LabelId {
	id := uuid.New()
	return &LabelId{g.UUIDWrapper[LabelId]{UUID: id}}

}

type Label struct {
	Id          LabelId
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ParentId    *LabelId
	Parent      *Label
}

func New(name, description string) (*Label, error) {
	id := NewLabelId()
	now := time.Now()
	label := &Label{Id: *id, Name: name,
		Description: description, CreatedAt: now,
		UpdatedAt: now}
	if err := label.Validate(); err != nil {
		return nil, err
	}
	return label, nil
}

func (l Label) Validate() error {
	if err := validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required),
		validation.Field(&l.Name, validation.Match(g.ResourceNameRegExp)),
	); err != nil {
		return fmt.Errorf("validating label name (%v): %w", l.Name, e.ErrResourceName)
	}
	return nil
}

func (l Label) String() string {
	labels := ""
	labels += l.Name
	if l.Parent != nil {
		labels += "," + l.Parent.String()
	}
	return labels
}
