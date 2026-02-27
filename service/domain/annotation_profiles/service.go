package annotation_profile

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	au "datahub/app/authorizer"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	clk "github.com/jonboulle/clockwork"
)

func NewAnnotationProfileService(annotationProfileRepo AnnotationProfileRepo,
	labelService *lbl.Service, logger *slog.Logger, auth *au.Authorizer,
	clock clk.Clock) *Service {
	return &Service{
		Repo:       annotationProfileRepo,
		Labels:     labelService,
		Logger:     logger,
		Authorizer: auth,
		Clock:      clock,
	}

}

type Service struct {
	Repo       AnnotationProfileRepo
	Labels     *lbl.Service
	Logger     *slog.Logger
	Authorizer *au.Authorizer
	Clock      clk.Clock
}

func (s *Service) Save(ctx context.Context, profile *AnnotationProfile) error {
	return s.Repo.Save(profile)
}

func (s *Service) ClearLabels(ctx context.Context, profile *AnnotationProfile) error {
	if err := s.Repo.ClearLabels(profile); err != nil {
		return err
	}
	profile.Labels = nil
	return nil
}

func (s *Service) GetLabelsOfProfile(ctx context.Context, profile *AnnotationProfile) ([]*lbl.Label, error) {
	labelIds, err := s.Repo.GetLabelIds(profile)
	if err != nil {
		return nil, err
	}

	var labels []*lbl.Label
	for _, id := range labelIds {
		label, err := s.Labels.Find(ctx, id)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	// Sort in place by Name
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})

	return labels, nil

}

func (s *Service) Find(ctx context.Context, id AnnotationProfileId) (*AnnotationProfile, error) {
	profile, err := s.Repo.Find(id)
	if err != nil {
		return nil, err
	}
	labels, err := s.GetLabelsOfProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	profile.Labels = labels
	return profile, nil

}
func (s *Service) List(ctx context.Context, pagination g.PaginationParams) ([]AnnotationProfile, *g.PaginationMeta, error) {

	profiles, meta, err := s.Repo.List(pagination)
	if err != nil {
		return nil, nil, err
	}
	for i, profile := range profiles {
		labels, err := s.GetLabelsOfProfile(ctx, &profile)
		if err != nil {
			return nil, nil, err
		}
		profiles[i].Labels = labels
	}
	return profiles, meta, nil
}

func (s *Service) Update(ctx context.Context, id AnnotationProfileId, payload ProfileUpdatables) (*AnnotationProfile, error) {

	baseErrMsg := "updating annotation profile"
	if err := s.Repo.Rename(id, payload.Name); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	profile, err := s.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return profile, nil

}

func (s *Service) FindByName(ctx context.Context, name string) (*AnnotationProfile, error) {
	profile, err := s.Repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	labels, err := s.GetLabelsOfProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	profile.Labels = labels
	return profile, nil

}

func (s *Service) findDuplicate(labelNames []string) *string {

	seen := make(map[string]bool)
	for _, v := range labelNames {
		if _, exists := seen[v]; exists {
			return &v
		}
		seen[v] = true
	}
	return nil
}

func (s *Service) AddLabel(ctx context.Context, profile *AnnotationProfile, label *lbl.Label) error {

	if profile.HasLabel(label) == true {
		return fmt.Errorf("adding label %v to profile %v: checking if label already exists: %w",
			label.Name, profile.Name,
			e.ErrDuplication)
	}

	profile.Labels = append(profile.Labels, label)
	return s.Repo.AddLabel(profile, label)
}

func (s *Service) AddLabelSet(ctx context.Context, profile *AnnotationProfile, labelNames []string) error {
	duplicate := s.findDuplicate(labelNames)
	if duplicate != nil {
		return fmt.Errorf("adding label set to profile %v: found duplicate label: %v: %w",
			profile.Name, *duplicate, e.ErrDuplication)
	}
	for _, name := range labelNames {
		label, err := s.Labels.FindByName(ctx, name)
		if err != nil {
			return err
		}
		err = s.AddLabel(ctx, profile, label)
		if err != nil {
			return fmt.Errorf("adding label set to profile %v: %w", profile.Name, err)
		}
	}
	return nil
}

func (s *Service) RemoveLabel(ctx context.Context, profile *AnnotationProfile, labelToRemove *lbl.Label) error {
	for i, label := range profile.Labels {
		if label.Id == labelToRemove.Id {
			profile.Labels = append(profile.Labels[:i], profile.Labels[i+1:]...)
			break
		}
	}
	return s.Repo.RemoveLabel(profile, labelToRemove)
}

func (s *Service) Delete(ctx context.Context, profile *AnnotationProfile) error {
	return s.Repo.Delete(profile)
}
