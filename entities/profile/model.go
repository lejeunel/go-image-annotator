package role

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type (
	ProfileName        = string
	ProfileDescription = string
)

type Profile struct {
	Id          ProfileId
	Name        ProfileName
	Description ProfileDescription
	Labels      []lbl.LabelName
}

func NewProfile(id ProfileId, name lbl.LabelName, opts ...Option) Profile {
	r := &Profile{Id: id, Name: name}
	for _, opt := range opts {
		opt(r)
	}
	return *r
}

type Option func(*Profile)

func WithDescription(d ProfileDescription) Option {
	return func(r *Profile) {
		r.Description = d
	}
}

func WithLabels(labels []lbl.LabelName) Option {
	return func(r *Profile) {
		r.Labels = labels
	}
}
