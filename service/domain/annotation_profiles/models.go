package annotation_profile

import (
	lbl "datahub/domain/labels"
	g "datahub/generic"
	"github.com/google/uuid"
)

type AnnotationProfileId struct {
	g.UUIDWrapper[AnnotationProfileId]
}

func NewAnnotationProfileId() *AnnotationProfileId {
	id := uuid.New()
	return &AnnotationProfileId{g.UUIDWrapper[AnnotationProfileId]{UUID: id}}

}

type AnnotationProfile struct {
	Id     AnnotationProfileId `db:"id"`
	Name   string              `db:"name"`
	Labels []*lbl.Label
}

func New(name string) *AnnotationProfile {
	return &AnnotationProfile{
		Id:   *NewAnnotationProfileId(),
		Name: name}

}

func (p *AnnotationProfile) LabelNames() []string {
	if p.Labels == nil {
		return []string{}
	}
	var names []string
	for _, label := range p.Labels {
		names = append(names, label.Name)
	}
	return names

}

func (p *AnnotationProfile) HasLabel(label *lbl.Label) bool {
	for _, l := range p.Labels {
		if l.Name == label.Name {
			return true
		}
	}
	return false

}

type ProfileUpdatables struct {
	Name string `json:"name" doc:"New profile name"`
}
