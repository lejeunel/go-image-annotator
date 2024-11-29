package services

import (
	"context"

	"github.com/google/uuid"
	e "go-image-annotator/errors"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
	"slices"
	"time"
)

type AnnotationService struct {
	LabelRepo       r.AnnotationRepo
	MaxPageSize     int
	DefaultPageSize int
}

func (s *AnnotationService) CreateLabel(ctx context.Context, label *m.Label) error {
	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return err
	}

	if err := label.Validate(); err != nil {
		return err
	}

	label.Id = uuid.New()

	label, err := s.LabelRepo.CreateLabel(ctx, label)
	if err != nil {
		return err
	}

	return nil

}

func (s *AnnotationService) GetLabelById(ctx context.Context, id string) (*m.Label, error) {

	label, err := s.LabelRepo.GetOneLabel(ctx, id)
	if err != nil {
		return nil, err
	}
	return label, nil

}

func (s *AnnotationService) DeleteLabel(ctx context.Context, label *m.Label) error {
	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return err
	}

	return s.LabelRepo.DeleteLabel(ctx, label)

}

func (s *AnnotationService) ApplyLabelToImage(ctx context.Context, label *m.Label, image *m.Image,
	collection *m.Collection) error {

	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return err
	}

	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return err
	}

	if err := s.LabelRepo.ApplyLabelToImage(ctx, label, image, collection, user.Email); err != nil {
		return err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image, collection)
	if err != nil {
		return err
	}

	image.Annotations = annotations

	return nil
}

func (s *AnnotationService) RemoveAnnotationFromImage(ctx context.Context, annotation *m.Annotation, image *m.Image, collection *m.Collection) error {
	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	isAuthorizedToDelete := (slices.Contains(user.Roles, "admin") || (user.Email == annotation.AuthorEmail))
	if !isAuthorizedToDelete {
		return e.ErrOwnershipPermission{Operation: "removing annotation from image", Details: "Only author of annotation can delete this."}
	}

	if err := s.LabelRepo.DeleteAnnotation(ctx, annotation); err != nil {
		return err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image, collection)
	if err != nil {
		return err
	}

	image.Annotations = annotations

	return nil
}

func (s *AnnotationService) ApplyBoundingBoxToImage(ctx context.Context, bbox *m.BoundingBox, image *m.Image, collection *m.Collection) error {
	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return err
	}

	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	if err := bbox.Validate(image); err != nil {
		return err
	}

	bbox.Id = uuid.New()
	bbox.ImageId = image.Id
	bbox.AuthorEmail = user.Email
	now := time.Now().String()
	bbox.CreatedAt = now
	bbox.UpdatedAt = now
	image.BoundingBoxes = append(image.BoundingBoxes, bbox)

	if err := s.LabelRepo.ApplyBoundingBoxToImage(ctx, bbox, image); err != nil {
		return err
	}
	return nil
}
