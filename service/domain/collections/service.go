package collections

import (
	"context"
	"slices"

	au "datahub/app/authorizer"
	pro "datahub/domain/annotation_profiles"
	lbl "datahub/domain/labels"
	g "datahub/generic"
	"encoding/json"
	"fmt"
	clk "github.com/jonboulle/clockwork"
	"log/slog"
)

type Service struct {
	CollectionRepo CollectionRepo
	Profiles       *pro.Service
	Labels         *lbl.Service
	MaxPageSize    int
	Authorizer     *au.Authorizer
	Clock          clk.Clock
	Logger         *slog.Logger
}

func NewCollectionService(repo CollectionRepo, profileService *pro.Service, labelService *lbl.Service, maxPageSize int, defaultPageSize int,
	auth *au.Authorizer, logger *slog.Logger, clock clk.Clock) *Service {
	return &Service{CollectionRepo: repo,
		Profiles:    profileService,
		Labels:      labelService,
		MaxPageSize: maxPageSize,
		Authorizer:  auth,
		Clock:       clock,
		Logger:      logger,
	}
}

func (s *Service) Create(ctx context.Context, collection *Collection) error {
	err := s.Authorizer.WantToContributeImages(ctx, collection.Group)

	if err != nil {
		return fmt.Errorf("creating collection: %w", err)
	}

	now := s.Clock.Now()
	collection.CreatedAt = now
	collection.UpdatedAt = now
	err = s.CollectionRepo.Create(collection)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) Touch(ctx context.Context, collection *Collection) error {
	collection.UpdatedAt = s.Clock.Now()
	err := s.CollectionRepo.Touch(collection.Id, s.Clock.Now())
	if err != nil {
		return fmt.Errorf("updating timestamp of collection: %w", err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, name string, collectionFields CollectionUpdatables) (*Collection, error) {

	_, err := New(collectionFields.Name, WithGroup(collectionFields.Group), WithDescription(collectionFields.Description))
	if err != nil {
		return nil, fmt.Errorf("validating collection fields: %w", err)
	}

	collection, err := s.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting collection %v prior to updating it: %w", collection, err)
	}

	err = s.Authorizer.WantToUpdateCollection(ctx, collection.Group)
	if err != nil {
		return nil, fmt.Errorf("checking authorization to update collection: %w", err)
	}

	err = s.CollectionRepo.Update(collection.Id, collectionFields)
	if err != nil {
		return nil, fmt.Errorf("updating collection: %w", err)
	}

	if err := s.Touch(ctx, collection); err != nil {
		return nil, fmt.Errorf("touching timestamp when updating collection: %w", err)
	}

	updatedCollection, err := s.Find(ctx, collection.Id)

	return updatedCollection, nil
}

func (s *Service) Patch(ctx context.Context, name string, patches g.JSONPatches) (*Collection, error) {
	original, err := s.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("patching collection: %w", err)
	}
	originalJSONBytes, err := json.Marshal(&original)
	if err != nil {
		return nil, fmt.Errorf("patching collection: %w", err)
	}

	modifiedBytes, err := patches.Apply(originalJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("patching collection: %w", err)
	}

	var modified CollectionUpdatables
	if err := json.Unmarshal(modifiedBytes, &modified); err != nil {
		return nil, fmt.Errorf("patching collection: %w", err)
	}

	return s.Update(ctx, name, modified)
}

func (s *Service) Find(ctx context.Context, id CollectionId) (*Collection, error) {
	baseErrMsg := "fetching collection by id"
	collection, err := s.CollectionRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("%v: %v: %w", baseErrMsg, id, err)
	}
	if err := s.appendProfile(ctx, collection); err != nil {
		return nil, fmt.Errorf("%v: %v: %w", baseErrMsg, id, err)
	}

	return collection, nil
}

func (s *Service) FindByName(ctx context.Context, name string) (*Collection, error) {
	baseErrMsg := "fetching collection by name"
	collection, err := s.CollectionRepo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%v: %v: %w", baseErrMsg, name, err)
	}
	if err := s.appendProfile(ctx, collection); err != nil {
		return nil, fmt.Errorf("%v: %v: %w", baseErrMsg, name, err)
	}
	return collection, nil
}

func (s *Service) Delete(ctx context.Context, id CollectionId) error {

	collection, err := s.CollectionRepo.Find(id)
	if err != nil {
		return fmt.Errorf("deleting collection: %w", err)
	}

	err = s.Authorizer.WantToContributeImages(ctx, collection.Group)
	if err != nil {
		return fmt.Errorf("deleting collection: %w", err)
	}
	return s.CollectionRepo.Delete(collection)
}

func (s *Service) List(
	ctx context.Context,
	ordering OrderingArgs,
	pagination g.PaginationParams) ([]Collection, *g.PaginationMeta, error) {

	if err := pagination.Validate(s.MaxPageSize); err != nil {
		return nil, nil, fmt.Errorf("listing collections: %w", err)
	}

	collections, paginationMeta, err := s.CollectionRepo.List(ordering, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("listing collections: %w", err)
	}

	for i := range len(collections) {
		err := s.appendProfile(ctx, &collections[i])
		if err != nil {
			return nil, nil, fmt.Errorf("listing collections: appending annotation profile: %w")
		}
	}

	return collections, paginationMeta, nil

}

func (s *Service) AssignProfile(ctx context.Context, profile *pro.AnnotationProfile, collection *Collection) error {
	collection.Profile = profile
	return s.CollectionRepo.AssignProfile(collection, profile)
}

func (s *Service) UnassignProfile(ctx context.Context, collection *Collection) error {
	collection.Profile = nil
	return s.CollectionRepo.UnassignProfile(collection)
}

func (s *Service) appendProfile(ctx context.Context, collection *Collection) error {
	if collection.ProfileId != nil {
		profile, err := s.Profiles.Find(ctx, *collection.ProfileId)
		if err != nil {
			return err
		}
		collection.Profile = profile
	}
	return nil

}

func (s *Service) GetAvailableLabels(ctx context.Context, collection *Collection) ([]lbl.Label, error) {
	if collection.Profile == nil {
		return s.Labels.GetAllLabels(ctx)
	}
	return g.MapDeref(collection.Profile.Labels), nil
}

func (s *Service) GetAvailableLabelNames(ctx context.Context, collection *Collection) ([]string, error) {
	labels, err := s.GetAvailableLabels(ctx, collection)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, l := range labels {
		names = append(names, l.Name)
	}

	return names, nil
}

func (s *Service) IsLabelAllowed(ctx context.Context, collection *Collection, label *lbl.Label) (bool, error) {
	if collection.Profile == nil {
		return true, nil
	}
	availableLabels, err := s.GetAvailableLabelNames(ctx, collection)
	if err != nil {
		return false, fmt.Errorf("checking whether label %v is applicable to collection %v: %w",
			label.Name, collection.Name, err)
	}

	if slices.Contains(availableLabels, label.Name) == false {
		return false, nil
	} else {
		return true, nil
	}

}
