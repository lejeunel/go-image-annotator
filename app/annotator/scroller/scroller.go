package scroller

import (
	"errors"
	"fmt"

	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

type Interface interface {
	Init(im.ImageId, ...Option) (*ScrollerState, error)
}

type Scroller struct {
	repo Repo
}

type ScrollerState struct {
	Next     *im.BaseImage
	Previous *im.BaseImage
}

func (s Scroller) Init(imageId im.ImageId, opts ...Option) (*ScrollerState, error) {
	criteria := NewCriteria(opts...)
	if err := checkCriteria(s.repo, imageId, criteria); err != nil {
		return nil, err
	}
	state := ScrollerState{}
	next, errNext := s.getOne(imageId, ScrollNext, criteria)
	prev, errPrev := s.getOne(imageId, ScrollPrevious, criteria)

	if errNext != nil || errPrev != nil {
		return nil, fmt.Errorf("%w, %w", errNext, errPrev)
	}
	state.Next = next
	state.Previous = prev

	return &state, nil
}

func (s *Scroller) getOne(current im.ImageId, direction ScrollingDirection, criteria ScrollingCriteria) (*im.BaseImage, error) {

	image, err := s.repo.GetAdjacent(current, criteria, direction)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		return nil, err
	}
	return image, nil
}

func checkCriteria(repo Repo, imageId im.ImageId, criteria ScrollingCriteria) error {
	errCtx := "initializing image scroller"
	if err := repo.ImageMustExist(imageId); err != nil {
		return fmt.Errorf("%v: checking that image with id %v exists: %w",
			errCtx, imageId, err)

	}
	if criteria.Collection != nil {
		if err := repo.CollectionMustExist(*criteria.Collection); err != nil {
			return fmt.Errorf("%v: checking that collection with name %v exists: %w",
				errCtx, *criteria.Collection, err)
		}
	}
	return nil
}

func New(repo Repo) Scroller {
	// criteria := NewCriteria(opts...)
	return Scroller{repo: repo}
}
