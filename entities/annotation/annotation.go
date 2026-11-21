package annotation

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type AnnotationId struct {
	uuidw.UUIDWrapper[AnnotationId]
}

func NewAnnotationId() AnnotationId {
	return AnnotationId{uuidw.UUIDWrapper[AnnotationId]{UUID: uuid.New()}}
}

func NewAnnotationIdFromString(s string) (*AnnotationId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid ImageId: %w", err)
	}

	return &AnnotationId{
		UUIDWrapper: uuidw.FromUUID[AnnotationId](id),
	}, nil
}

type Points struct {
	Coordinates [][2]float32
}

func (p *Points) Append(x, y float32) {
	p.Coordinates = append(p.Coordinates, [2]float32{x, y})
}

func (p Points) MaxX() float32 {
	if len(p.Coordinates) == 0 {
		return float32(math.Inf(1))
	}

	r := float32(math.Inf(-1))
	for _, p := range p.Coordinates {
		if p[0] > r {
			r = p[0]
		}
	}
	return r
}

func (p Points) MaxY() float32 {
	if len(p.Coordinates) == 0 {
		return float32(math.Inf(1))
	}

	r := float32(math.Inf(-1))
	for _, p := range p.Coordinates {
		if p[1] > r {
			r = p[1]
		}
	}
	return r
}

func (p Points) MinX() float32 {
	if len(p.Coordinates) == 0 {
		return float32(math.Inf(-1))
	}

	r := float32(math.Inf(1))
	for _, p := range p.Coordinates {
		if p[0] < r {
			r = p[0]
		}
	}
	return r
}

func (p Points) MinY() float32 {
	if len(p.Coordinates) == 0 {
		return float32(math.Inf(-1))
	}

	r := float32(math.Inf(1))
	for _, p := range p.Coordinates {
		if p[1] < r {
			r = p[1]
		}
	}
	return r
}

type ImageLabel struct {
	Id     AnnotationId
	Label  lbl.Label
	Author *u.UserId
	Time   *time.Time
}

type Annotation struct {
	Id    AnnotationId
	Label string
}

type BoundingBox struct {
	Id     AnnotationId
	Label  lbl.Label
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
	Angle  float32
	Author *u.UserId
	Time   *time.Time
}

type Polygon struct {
	Id     AnnotationId
	Label  lbl.Label
	Points Points
	Author *u.UserId
	Time   *time.Time
}

type PolygonRequest struct {
	Label  string
	Points Points
}

type PolygonUpdatables struct {
	LabelId lbl.LabelId
	Points  Points
}

type BoundingBoxResponse struct {
	Label  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

type BoundingBoxRequest struct {
	Label  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

type BoundingBoxUpdatables struct {
	LabelId lbl.LabelId
	Xc      float32
	Yc      float32
	Width   float32
	Height  float32
	Angle   float32
}

type Option func(*BoundingBox)

func WithAngle(ang float32) Option {
	return func(c *BoundingBox) {
		c.Angle = ang
	}
}

func NewBoundingBox(id AnnotationId, xc float32, yc float32,
	width float32, height float32, label lbl.Label, opts ...Option,
) BoundingBox {
	b := &BoundingBox{Id: id, Xc: xc, Yc: yc, Width: width, Height: height, Label: label}
	for _, opt := range opts {
		opt(b)
	}
	return *b
}

func ValidateBoundingBox(
	xc float32,
	yc float32,
	width float32,
	height float32,
	angle float32,
) error {
	errCtx := "validating bounding box"
	if width <= 0 {
		return fmt.Errorf(
			"%v: checking whether width (%v) <= 0: %w",
			errCtx,
			width,
			e.ErrValidation,
		)
	}
	return nil
}

// minPolygonArea is the smallest area, in squared pixels, that a polygon
// may enclose. Comparing with a tolerance rather than with zero absorbs
// floating-point rounding: collinear points with non-integer coordinates
// can yield a tiny, non-zero area.
const minPolygonArea = 1e-6

// ValidatePolygon checks that points describe a usable polygon: at least
// three points, finite coordinates, and a non-zero enclosed area (see
// Points.Area).
func ValidatePolygon(points Points) error {
	errCtx := "validating polygon"
	n := len(points.Coordinates)
	if n < 3 {
		return fmt.Errorf(
			"%v: checking whether number of points (%v) < 3: %w",
			errCtx,
			n,
			e.ErrValidation,
		)
	}
	for i, p := range points.Coordinates {
		for _, c := range p {
			if math.IsNaN(float64(c)) || math.IsInf(float64(c), 0) {
				return fmt.Errorf(
					"%v: checking whether point %v (%v) has finite coordinates: %w",
					errCtx,
					i,
					p,
					e.ErrValidation,
				)
			}
		}
	}
	if area := points.Area(); area < minPolygonArea {
		return fmt.Errorf(
			"%v: checking whether area (%v) is zero: points are collinear or repeated: %w",
			errCtx,
			area,
			e.ErrValidation,
		)
	}
	return nil
}

// Area returns the area enclosed by the polygon whose vertices are the
// points, taken in order, using the shoelace formula (also called Gauss's
// area formula)
func (p Points) Area() float64 {
	n := len(p.Coordinates)
	var sum float64
	for i := range n {
		cur, next := p.Coordinates[i], p.Coordinates[(i+1)%n]
		sum += float64(cur[0])*float64(next[1]) - float64(next[0])*float64(cur[1])
	}
	return math.Abs(sum) / 2
}

func NewImageLabel(label lbl.Label) ImageLabel {
	return ImageLabel{Id: NewAnnotationId(), Label: label}
}

func NewPolygon(id AnnotationId, points Points, label lbl.Label) Polygon {
	return Polygon{Id: id, Points: points, Label: label}
}
