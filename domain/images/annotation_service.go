package images

import (
	"context"
	"fmt"
	"log/slog"

	au "datahub/app/authorizer"
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
	e "datahub/errors"

	"github.com/google/uuid"
	clk "github.com/jonboulle/clockwork"
)

func NewAnnotationService(annotationRepo AnnotationRepo,
	collectionService *clc.Service,
	labelService *lbl.Service, logger *slog.Logger, auth *au.Authorizer,
	clock clk.Clock) *AnnotationService {
	return &AnnotationService{
		AnnotationRepo: annotationRepo,
		Collections:    collectionService,
		Labels:         labelService,
		Logger:         logger,
		Authorizer:     auth,
		Clock:          clock,
	}

}

type AnnotationService struct {
	AnnotationRepo AnnotationRepo
	Collections    *clc.Service
	Labels         *lbl.Service
	Logger         *slog.Logger
	Authorizer     *au.Authorizer
	Clock          clk.Clock
}

func (s *AnnotationService) GetAnnotations(ctx context.Context, image *Image) ([]*Annotation, error) {
	annotationsIds, err := s.AnnotationRepo.GetAnnotationIdsOfImage(image)
	if err != nil {
		return nil, fmt.Errorf("appending annotations to image %v of collection %v: fetching annotations: %w",
			image.Id, image.CollectionId, err)
	}

	var annotations []*Annotation

	for _, id := range annotationsIds {
		a, err := s.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if a.LabelId.UUID != uuid.Nil {
			label, err := s.Labels.Find(ctx, a.LabelId)
			if err != nil {
				return nil, err
			}
			a.Label = label

		}
		annotations = append(annotations, a)

	}
	return annotations, nil
}

func (s *AnnotationService) AppendAnnotations(ctx context.Context, image *Image) error {
	annotations, err := s.GetAnnotations(ctx, image)
	if err != nil {
		return err
	}
	image.Annotations = annotations

	bboxes, err := s.AnnotationRepo.GetBoundingBoxesOfImage(image)
	for _, b := range bboxes {
		label, err := s.Labels.Find(ctx, b.Annotation.LabelId)
		if err != nil {
			return fmt.Errorf("appending label to bounding-box: %w", err)
		}
		b.Annotation.Label = label
	}
	if err != nil {
		return fmt.Errorf("appending annotations to image %v of collection %v: fetching bounding boxes: %w",
			image.Id, image.CollectionId, err)
	}
	image.BoundingBoxes = bboxes

	return nil

}

func (s *AnnotationService) UpdateLabel(ctx context.Context, annotationId string, labelName string) error {
	label, err := s.Labels.FindByName(ctx, labelName)
	if err != nil {
		return fmt.Errorf("updating annotation label: %w", err)
	}
	return s.AnnotationRepo.UpdateAnnotationLabel(annotationId, label.Id.String())
}

func (s *AnnotationService) ApplyLabel(ctx context.Context, label *lbl.Label, image *Image) error {

	if err := s.Authorizer.WantToContributeAnnotations(ctx, image.Group); err != nil {
		return fmt.Errorf("applying label: %w", err)
	}

	email, err := s.Authorizer.IdentityProvider.Email(ctx)
	if err != nil {
		return fmt.Errorf("fetching email from context: %w", err)
	}

	if err := s.AnnotationRepo.ApplyLabelToImage(label, image, email); err != nil {
		return fmt.Errorf("applying label to image: %w", err)
	}

	annotations, err := s.GetAnnotations(ctx, image)
	if err != nil {
		return err
	}

	image.Annotations = annotations

	return nil
}

func (s *AnnotationService) Find(ctx context.Context, id string) (*Annotation, error) {
	return s.AnnotationRepo.GetAnnotationById(id)
}

func (s *AnnotationService) Delete(ctx context.Context, id string) error {
	annotation, err := s.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching annotation prior to deleting it (id: %v): %w", id, err)
	}
	if err := s.AnnotationRepo.DeleteAnnotation(annotation); err != nil {
		return fmt.Errorf("deleting annotation by id (%v): %w", id, err)
	}
	return nil

}

func (s *AnnotationService) applyBoundingBox(ctx context.Context, bbox *BoundingBox, image *Image) error {

	image.BoundingBoxes = append(image.BoundingBoxes, bbox)

	if err := s.AnnotationRepo.ApplyBoundingBox(bbox, image); err != nil {
		return err
	}

	collection, err := s.Collections.Find(ctx, image.CollectionId)
	if err != nil {
		return fmt.Errorf("apply bounding box: fetching collection: %w", err)
	}

	if err := s.Collections.Touch(ctx, collection); err != nil {
		return fmt.Errorf("apply bounding box: touching collection: %w", err)
	}

	return nil
}
func (s *AnnotationService) imageHasBoundingBoxWithId(image *Image, bboxId uuid.UUID) bool {
	for _, b := range image.BoundingBoxes {
		if b.Annotation.Id == bboxId {
			return true
		}
	}
	return false
}
func (s *AnnotationService) updateBoundingBox(ctx context.Context, bbox *BoundingBox, image *Image) error {
	baseErrMsg := "updating bounding box"

	bbox.Annotation.UpdatedAt = s.Clock.Now()
	if err := s.AnnotationRepo.UpdateBoundingBox(bbox, image); err != nil {
		return fmt.Errorf("%v updating bounding-box: %w", baseErrMsg, err)
	}
	collection, err := s.Collections.Find(ctx, image.CollectionId)
	if err != nil {
		return fmt.Errorf("%v: fetching collection: %w", baseErrMsg, err)
	}
	if err := s.Collections.Touch(ctx, collection); err != nil {
		return fmt.Errorf("%v: touching collection: %w", baseErrMsg, err)
	}
	return nil
}
func (s *AnnotationService) UpsertBoundingBox(ctx context.Context, bbox *BoundingBox, image *Image) error {
	errCtx := fmt.Sprintf("upserting bounding box %v on image %v", bbox.Annotation.Id,
		image.Id)
	if err := s.Authorizer.WantToContributeAnnotations(ctx, image.Group); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	if bbox.Annotation.Label == nil {
		return fmt.Errorf("%v: has no assigned label: %w", errCtx, e.ErrValidation)
	}

	collection, err := s.Collections.Find(ctx, image.CollectionId)
	if err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	isLabelAllowed, err := s.Collections.IsLabelAllowed(ctx, collection, bbox.Annotation.Label)
	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, "checking for applicable labels", err)
	}

	if isLabelAllowed == false {
		return fmt.Errorf("%v: %v: %w", errCtx, "checking for applicable labels", e.ErrForbiddenLabel)
	}

	email, err := s.Authorizer.IdentityProvider.Email(ctx)
	if err != nil {
		return fmt.Errorf("%v: fetching identity: %w", errCtx, err)
	}

	bbox.Annotation.ImageId = image.Id
	bbox.Annotation.AuthorEmail = email
	bbox.Annotation.CollectionId = image.CollectionId

	if s.imageHasBoundingBoxWithId(image, bbox.Annotation.Id) {
		if err := s.updateBoundingBox(ctx, bbox, image); err != nil {
			return fmt.Errorf("%v: %w", errCtx, err)
		}
		return nil
	}

	bbox.Annotation.CreatedAt = s.Clock.Now()
	bbox.Annotation.UpdatedAt = s.Clock.Now()
	if err := s.applyBoundingBox(ctx, bbox, image); err != nil {
		return fmt.Errorf("%v: inserting bounding-box: %w", errCtx, err)
	}

	return nil
}
