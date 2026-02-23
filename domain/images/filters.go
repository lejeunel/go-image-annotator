package images

import (
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strings"
)

type FilterArgs struct {
	CollectionId   *clc.CollectionId `json:"collection_id" doc:"Id of collection"`
	CollectionName *string           `json:"collection_name" doc:"Name of collection"`
	CameraId       *loc.CameraId     `json:"camera_id" doc:"Id of camera"`
	LabelId        *lbl.LabelId      `json:"label_id" doc:"Id of label"`
	TemporalFilter *TemporalFilter
}

type TemporalFilter struct {
	ReferenceImageId ImageId
	Field            string
	Before           bool
}

func (f *FilterArgs) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	q = q.PlaceholderFormat(sq.Question)
	if f.CollectionId != nil {
		q = q.Where(sq.Eq{"ic.collection_id": *f.CollectionId})
	}

	if f.CollectionName != nil {
		subq := sq.StatementBuilder.Select("id").From("collections").Where(sq.Eq{"name": *f.CollectionName})
		q = q.Where(sq.Expr("ic.collection_id=(?)", subq))
	}

	if f.CameraId != nil {
		q = q.Where(sq.Eq{"i.camera_id": *f.CameraId})
	}

	if f.LabelId != nil {
		subq := sq.StatementBuilder.Select("image_id").From("annotations").Where(sq.Eq{"label_id": *f.LabelId})
		q = q.Where(sq.Expr("i.id IN (?)", subq))
	}

	if f.TemporalFilter != nil {
		subq := sq.StatementBuilder.Select(f.TemporalFilter.Field).From("images").Where(sq.Eq{"id": f.TemporalFilter.ReferenceImageId})
		if f.TemporalFilter.Before {
			q = q.Where(sq.Expr(fmt.Sprintf("i.%v < (?)", f.TemporalFilter.Field), subq))
			q = q.OrderBy(fmt.Sprintf("i.%v DESC", f.TemporalFilter.Field))
		} else {
			q = q.Where(sq.Expr(fmt.Sprintf("i.%v > (?)", f.TemporalFilter.Field), subq))
			q = q.OrderBy(fmt.Sprintf("i.%v ASC", f.TemporalFilter.Field))
		}
	}

	q = q.PlaceholderFormat(sq.Dollar)
	return q

}

func (f *FilterArgs) String() string {
	res := []string{}
	if f.CollectionId != nil {
		res = append(res, fmt.Sprintf("collection id: %v", *f.CollectionId))
	}
	if f.CollectionName != nil {
		res = append(res, fmt.Sprintf("collection name: %v", *f.CollectionName))
	}
	if f.LabelId != nil {
		res = append(res, fmt.Sprintf("label id: %v", *f.LabelId))
	}
	if f.CameraId != nil {
		res = append(res, fmt.Sprintf("camera id: %v", *f.CameraId))
	}

	return strings.Join(res, " / ")
}

func NewImageFilterFromString(collectionId, collectionName, cameraId, labelId string) (*FilterArgs, error) {
	var opts []FilterOption
	if labelId != "" {
		labelId, err := lbl.NewLabelIdFromString(labelId)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithLabelId(*labelId))
	}
	if collectionId != "" {
		collectionId, err := clc.NewCollectionIdFromString(collectionId)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithCollectionId(*collectionId))
	}

	opts = append(opts, WithCollectionName(collectionName))

	if cameraId != "" {
		cameraId, err := loc.NewCameraIdFromString(cameraId)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithCameraId(*cameraId))
	}
	return NewImageFilter(opts...), nil

}

type FilterOption func(*FilterArgs)

func WithCollectionId(id clc.CollectionId) FilterOption {
	return func(f *FilterArgs) {
		f.CollectionId = &id
	}
}

func WithCollectionName(name string) FilterOption {
	return func(f *FilterArgs) {
		if name != "" {
			f.CollectionName = &name
		}
	}
}

func WithCameraId(id loc.CameraId) FilterOption {
	return func(f *FilterArgs) {
		f.CameraId = &id
	}
}

func WithLabelId(id lbl.LabelId) FilterOption {
	return func(f *FilterArgs) {
		f.LabelId = &id
	}
}

func NewImageFilter(opts ...FilterOption) *FilterArgs {
	filter := &FilterArgs{}
	for _, opt := range opts {
		opt(filter)
	}
	return filter
}

type OrderingArgs struct {
	CapturedAt *bool
	CreatedAt  *bool
	Descending bool
}

func NewImageDefaultOrderingArgs() *OrderingArgs {
	captured_at := true
	return &OrderingArgs{CapturedAt: &captured_at}
}

func (o *OrderingArgs) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	if o.CapturedAt != nil {
		if o.Descending == true {
			q = q.OrderBy("i.captured_at DESC")

		} else {
			q = q.OrderBy("i.captured_at ASC")

		}
	}
	if o.CreatedAt != nil {
		if o.Descending == true {
			q = q.OrderBy("i.created_at DESC")

		} else {
			q = q.OrderBy("i.created_at ASC")

		}
	}

	// this is necessary to impose consistency, without an explicit order by,
	// the DB engine could return two different ordering given the same query.
	q = q.OrderBy("i.id ASC")

	return q
}

func NewAscendingImageCapturedOrder() *OrderingArgs {
	c := true
	return &OrderingArgs{CapturedAt: &c, Descending: false}
}

func NewDescendingImageCapturedOrder() *OrderingArgs {
	c := true
	return &OrderingArgs{CapturedAt: &c, Descending: true}
}

func NewAscendingImageCreatedOrder() *OrderingArgs {
	c := true
	return &OrderingArgs{CreatedAt: &c, Descending: false}
}

func NewDescendingImageCreatedOrder() *OrderingArgs {
	c := true
	return &OrderingArgs{CreatedAt: &c, Descending: true}
}
