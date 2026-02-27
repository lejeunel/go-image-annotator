package annotator

import (
	"context"
	im "datahub/domain/images"
	e "datahub/errors"
)

type Scroller struct {
	Images   *im.Service
	Filters  im.FilterArgs
	Ordering im.OrderingArgs
}

func NewScroller(i *im.Service, filters im.FilterArgs, ordering im.OrderingArgs) *Scroller {
	return &Scroller{Images: i, Filters: filters, Ordering: ordering}
}

func (s *Scroller) GetNextImage(ctx context.Context, current *im.Image) (*im.Image, error) {
	next, err := s.Images.GetAdjacent(ctx, current, s.Filters,
		s.Ordering, false, im.FetchMetaOnly)
	if err != nil {
		return nil, err
	}
	if next == nil {
		return nil, e.ErrNotFound
	}
	return next, nil
}

func (s *Scroller) GetPrevImage(ctx context.Context, current *im.Image) (*im.Image, error) {
	prev, err := s.Images.GetAdjacent(ctx, current, s.Filters,
		s.Ordering, true, im.FetchMetaOnly)
	if err != nil {
		return nil, err
	}
	if prev == nil {
		return nil, e.ErrNotFound
	}
	return prev, nil
}
