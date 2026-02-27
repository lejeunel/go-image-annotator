package annotator

import (
	"fmt"
	"golang.org/x/image/draw"
	goim "image"
	jpeg "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

type Rescaler struct {
	TargetWidth int
}

func (t *Rescaler) TransformImage(r io.Reader, w io.Writer) error {

	src, _, err := goim.Decode(r)
	if err != nil {
		return fmt.Errorf("applying transformation to image: %w", err)
	}
	ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	height := int(math.Round(float64(t.TargetWidth) * ratio))

	dst := goim.NewRGBA(goim.Rect(0, 0, t.TargetWidth, height))
	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	jpeg.Encode(w, dst, nil)
	return nil
}

func (t *Rescaler) ForwardTransformBoundingBoxCoords(origWidth float64, bbox *BoundingBox) error {
	factor := float64(t.TargetWidth) / float64(origWidth)
	bbox.Xc = bbox.Xc * factor
	bbox.Yc = bbox.Yc * factor
	bbox.Height = bbox.Height * factor
	bbox.Width = bbox.Width * factor

	return nil
}

func (t *Rescaler) BackwardTransformBoundingBoxCoords(origWidth float64, bbox *BoundingBox) error {
	factor := float64(t.TargetWidth) / float64(origWidth)
	bbox.Xc = bbox.Xc / factor
	bbox.Yc = bbox.Yc / factor
	bbox.Height = bbox.Height / factor
	bbox.Width = bbox.Width / factor

	return nil
}
