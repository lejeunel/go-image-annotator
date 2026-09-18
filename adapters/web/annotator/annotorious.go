package annotator

import (
	"bytes"
	"text/template"

	_ "embed"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed templates/annotator.js
var AnnotatorJS string

type AnnotatorData struct {
	URLs             AnnotatorURLs
	ImageId          string
	Collection       string
	EnableAnnotation bool
}

type AnnotatorURLs = struct {
	FetchAnnotations string
	SetLabel         string
	SubmitImageLabel string
	SubmitBox        string
	SubmitPolygon    string
	UpdateBox        string
	UpdatePolygon    string
	RemoveAnnotation string
	AnnotationPanel  string
}

var annotatorURLs = AnnotatorURLs{
	Annotations,
	SetLabel,
	SubmitImageLabel,
	SubmitBox,
	SubmitPolygon,
	UpdateBox,
	UpdatePolygon,
	RemoveAnnotation,
	AnnotationPanel,
}

func MakeAnnotoriousScript(imageId string, collection string) (*Node, error) {
	tAnnot, err := template.New("annotator").Parse(AnnotatorJS)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString("")
	data := AnnotatorData{
		URLs:             annotatorURLs,
		ImageId:          imageId,
		Collection:       collection,
		EnableAnnotation: true,
	}

	err = tAnnot.ExecuteTemplate(buf, "annotator", data)
	if err != nil {
		return nil, err
	}
	script := Script(Raw(buf.String()), Defer())
	return &script, nil
}
