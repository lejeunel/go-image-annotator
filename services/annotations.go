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

func (s *AnnotationService) Create(ctx context.Context, label *m.Label) (*m.Label, error) {
	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return nil, err
	}

	if err := label.Validate(); err != nil {
		return nil, err
	}

	label.Id = uuid.New()

	label, err := s.LabelRepo.CreateLabel(ctx, label)
	if err != nil {
		return nil, err
	}

	return label, nil

}

func (s *AnnotationService) GetOne(ctx context.Context, id string) (*m.Label, error) {

	label, err := s.LabelRepo.GetOneLabel(ctx, id)
	if err != nil {
		return nil, err
	}
	return label, nil

}

func (s *AnnotationService) Delete(ctx context.Context, label *m.Label) error {

	numImages, err := s.LabelRepo.NumImagesWithLabel(ctx, label)
	if err != nil {
		return err
	}

	if numImages > 0 {
		return e.ErrForbiddenDeletingDependency{ParentEntity: "label", ParentId: label.Id.String(), ChildEntity: "image"}
	}

	return s.LabelRepo.DeleteLabel(ctx, label)

}

func (s *AnnotationService) ApplyLabelToImage(ctx context.Context, label *m.Label, image *m.Image) (*m.Image, error) {

	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return nil, err
	}

	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return image, err
	}

	if err := s.LabelRepo.ApplyLabelToImage(ctx, label, image, user.Email); err != nil {
		return nil, err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}

	image.Annotations = annotations

	return image, nil
}

func (s *AnnotationService) RemoveAnnotationFromImage(ctx context.Context, annotation *m.ImageAnnotation, image *m.Image) (*m.Image, error) {
	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return image, err
	}
	isAuthorizedToDelete := (slices.Contains(user.Roles, "admin") || (user.Email == annotation.AuthorEmail))
	if !isAuthorizedToDelete {
		return image, e.ErrOwnershipPermission{Operation: "removing annotation from image", Details: "Only author of annotation can delete this."}
	}

	if err := s.LabelRepo.RemoveAnnotationFromImage(ctx, annotation); err != nil {
		return image, err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image)
	if err != nil {
		return image, err
	}

	image.Annotations = annotations

	return image, nil
}

func (s *AnnotationService) ApplyPolygonToImage(ctx context.Context, polygon *m.Polygon, image *m.Image) (*m.Image, error) {
	if err := g.CheckAuthorization(ctx, "annotation-contrib"); err != nil {
		return image, err
	}

	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return image, err
	}
	polygon.AuthorEmail = user.Email
	now := time.Now().String()
	polygon.CreatedAt = now
	polygon.UpdatedAt = now
	image.Polygons = append(image.Polygons, polygon)

	if err := s.LabelRepo.ApplyPolygonToImage(ctx, polygon, image); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *AnnotationService) DeletePolygonFromImage(ctx context.Context, polygon *m.Polygon, image *m.Image) (*m.Image, error) {

	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return image, err
	}
	isAuthorizedToDelete := (slices.Contains(user.Roles, "admin") || (user.Email == polygon.AuthorEmail))
	if !isAuthorizedToDelete {
		return image, e.ErrOwnershipPermission{Operation: "removing polygon from image", Details: "Only author of polygon can delete this."}
	}

	if err := s.LabelRepo.DeletePolygon(ctx, polygon); err != nil {
		return image, err
	}

	polygons, err := s.LabelRepo.GetPolygonsOfImage(ctx, image)
	if err != nil {
		return image, err
	}

	image.Polygons = polygons

	return image, nil
}
