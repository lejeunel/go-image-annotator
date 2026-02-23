package labels

import (
	"context"
	"fmt"
	"log/slog"

	au "datahub/app/authorizer"
	g "datahub/generic"
	clk "github.com/jonboulle/clockwork"
)

type Service struct {
	Repo            Repo
	MaxPageSize     int
	DefaultPageSize int
	Authorizer      *au.Authorizer
	Logger          *slog.Logger
	Clock           clk.Clock
}

type Updatables struct {
	Description string `json:"description" doc:"Description"`
}

func NewLabelService(repo Repo, maxPageSize int, defaultPageSize int,
	auth *au.Authorizer, logger *slog.Logger, clock clk.Clock) *Service {
	return &Service{Repo: repo,
		MaxPageSize:     maxPageSize,
		DefaultPageSize: defaultPageSize,
		Authorizer:      auth,
		Clock:           clock,
		Logger:          logger,
	}
}

func (s *Service) Create(ctx context.Context, label *Label) error {

	baseErrMsg := fmt.Sprintf("creating label with name %v", label.Name)
	err := s.Authorizer.WantToContributeLabels(ctx)
	if err != nil {
		err = fmt.Errorf("%v: %w", baseErrMsg, err)
		s.Logger.Error(err.Error())
		return err
	}

	err = s.Repo.Create(label)
	if err != nil {
		return fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return nil

}

func (s *Service) Delete(ctx context.Context, label *Label) error {
	err := s.Authorizer.WantToContributeLabels(ctx)
	if err != nil {
		err = fmt.Errorf("deleting label: %w", err)
		s.Logger.Error(err.Error())
		return err
	}

	if err := s.Repo.Delete(label); err != nil {
		err = fmt.Errorf("deleting label %v: %w", label.Name, err)
		s.Logger.Error(err.Error())
		return err

	}
	return nil

}

func (s *Service) Update(ctx context.Context, label *Label, values Updatables) error {
	if err := s.Repo.Update(label, values); err != nil {
		err = fmt.Errorf("updating label %v: %w", label.Name, err)
		s.Logger.Error(err.Error())
		return err
	}
	label.Description = values.Description
	return nil
}

func (s *Service) FindByName(ctx context.Context, name string) (*Label, error) {
	label, err := s.Repo.FindByName(name)
	if err != nil {
		err = fmt.Errorf("fetching label by name %v: %w", name, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	return label, nil
}

func (s *Service) Find(ctx context.Context, id LabelId) (*Label, error) {
	label, err := s.Repo.Find(id)
	if err != nil {
		err = fmt.Errorf("getting label by id (%v): %w", id, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	if err := s.populateParentLabel(ctx, label); err != nil {
		err = fmt.Errorf("populating parenting of label %v: %w", label.Name, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	return label, nil
}

func (s *Service) GetAllLabels(ctx context.Context) ([]Label, error) {

	numLabels, err := s.Repo.Count()
	if err != nil {
		err = fmt.Errorf("fetching all label names: %w", err)
		s.Logger.Error(err.Error())
		return nil, err
	}

	labels, _, err := s.List(ctx, g.OrderingArg{Field: "name"},
		g.PaginationParams{Page: 1, PageSize: int(numLabels)})
	if err != nil {
		err = fmt.Errorf("fetching all label names: %w", err)
		s.Logger.Error(err.Error())
		return nil, err
	}

	return labels, nil
}

func (s *Service) populateParentLabel(ctx context.Context, label *Label) error {

	if label.ParentId != nil {

		parent, err := s.Repo.Find(*NewLabelIdFromUUID(label.ParentId.UUID))
		if err != nil {
			return fmt.Errorf("fetching label: %w", err)
		}
		label.Parent = parent
		if err = s.populateParentLabel(ctx, parent); err != nil {
			return fmt.Errorf("fetching label: %w", err)
		}
	}

	return nil

}

func (s *Service) List(
	ctx context.Context,
	ordering g.OrderingArg,
	pagination g.PaginationParams) ([]Label, *g.PaginationMeta, error) {

	return s.Repo.List(ordering, pagination)
}

func (s *Service) SetParenting(ctx context.Context, child *Label, parent *Label) error {
	child.Parent = parent
	child.ParentId = &parent.Id
	if err := s.Repo.SetParenting(child, parent); err != nil {
		return fmt.Errorf("setting parenting relationship. Child id: %v, parent id: %v: %w", child.Id, parent.Id, err)
	}

	return nil
}
