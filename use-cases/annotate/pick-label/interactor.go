package pick

import (
	"context"
	"fmt"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

var defaultLabelCountLimit = 200

type Interface interface {
	Execute(context.Context, *pr.ProfileName, OutputPort)
}
type Interactor struct {
	LabelRepo
	ProfileRepo
	countLimit int
}

func (i Interactor) Execute(ctx context.Context, profileName *pr.ProfileName, out OutputPort) {
	errCtx := fmt.Errorf("listing labels")
	if profileName == nil {
		labels, err := i.fetchAll()
		if err != nil {
			out.Error(fmt.Errorf("%w: %w", errCtx, err))
			return
		}
		out.SuccessFetchLabels(labels)
		return
	}

	profile, err := i.ProfileRepo.Find(*profileName)
	if err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}
	out.SuccessFetchLabels(profile.Labels)
}

func (i Interactor) fetchAll() ([]string, error) {
	count, err := i.LabelRepo.Count()
	if err != nil {
		return nil, err
	}
	if count > int64(i.countLimit) {
		return nil, fmt.Errorf("checking whether current label count (%v) exceeds limit (%v): %w",
			count, i.countLimit, e.ErrLabelLimitExceeded)
	}

	labels, err := i.LabelRepo.FetchAll()
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func New(lr LabelRepo, pr ProfileRepo, opts ...Option) Interactor {
	i := &Interactor{
		LabelRepo:   lr,
		ProfileRepo: pr,
		countLimit:  defaultLabelCountLimit,
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

type Option func(*Interactor)

func WithLimit(limit int) Option {
	return func(c *Interactor) {
		c.countLimit = limit
	}
}
