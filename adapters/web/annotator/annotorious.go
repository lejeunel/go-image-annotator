package annotator

import (
	"bytes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"text/template"
)

type Bounds struct {
	MinX float32 `json:"minX"`
	MinY float32 `json:"minY"`
	MaxX float32 `json:"maxX"`
	MaxY float32 `json:"maxY"`
}

type BoxGeometry struct {
	XTopLeft float32 `json:"x"`
	YTopLeft float32 `json:"y"`
	W        float32 `json:"w"`
	H        float32 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
	Rot      float32 `json:"rot"`
}

type BoxSelector struct {
	Type     string      `json:"type"`
	Geometry BoxGeometry `json:"geometry"`
}

type BoxTarget struct {
	Selector BoxSelector `json:"selector"`
}

type AnnotoriousBody struct {
	Purpose string `json:"purpose"`
	Value   string `json:"value"`
}

type AnnotoriousBoxModel struct {
	AnnotationId string            `json:"id"`
	Properties   Properties        `json:"properties"`
	Target       BoxTarget         `json:"target"`
	Bodies       []AnnotoriousBody `json:"bodies"`
}

func (b AnnotoriousBoxModel) ExtractCoordinates() BoxCoordinates {
	return BoxCoordinates{
		Xc:     b.Target.Selector.Geometry.XTopLeft + b.Target.Selector.Geometry.W/2,
		Yc:     b.Target.Selector.Geometry.YTopLeft + b.Target.Selector.Geometry.H/2,
		Width:  b.Target.Selector.Geometry.W,
		Height: b.Target.Selector.Geometry.H,
		Angle:  b.Target.Selector.Geometry.Rot,
	}
}

type Properties struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type AnnotoriousBoxRequest struct {
	ImageId    string              `json:"image_id"`
	Collection string              `json:"collection"`
	Annotation AnnotoriousBoxModel `json:"annotation"`
	Label      string              `json:"label"`
}

type BoxCoordinates struct {
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
	Angle  float32
}

type AnnotatorState struct {
	ImageId          string
	Collection       string
	EnableAnnotation bool
}

func MakeAnnotoriousScript(imageId string, collection string) (*Node, error) {
	tAnnot, err := template.New("annotator").ParseFS(templatesFiles, "templates/annotator.js")
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString("")
	data := AnnotatorState{ImageId: imageId,
		Collection:       collection,
		EnableAnnotation: true}

	err = tAnnot.ExecuteTemplate(buf, "annotator", data)
	if err != nil {
		return nil, err
	}
	script := Script(Raw(buf.String()))
	return &script, nil
}
