package annotator

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	rt "github.com/lejeunel/go-image-annotator/routes"
	addmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/add"
	delmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/delete"
	listmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/list"
	readmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/read"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var (
	MetaDivId         = "meta"
	MetaFormId        = "create-meta"
	MetaTypeQueryArg  = "type"
	MetaTypeDateTime  = "datetime"
	MetaTypeInt       = "int"
	MetaTypeFloat     = "float"
	MetaTypeBoolean   = "boolean"
	MetaTypeString    = "string"
	MetaImageIdArg    = "image_id"
	MetaCollectionArg = "collection"
	MetaKeyArg        = "key"
	MetaValueArg      = "value"
	MetaTableFields   = []string{"key", "value", ""}
)

func parseFormValue(metaType, value string) any {
	var parsed any
	var err error
	switch metaType {
	case MetaTypeBoolean:
		if value == "on" {
			return true
		} else {
			return false
		}
	case MetaTypeInt:
		parsed, err = strconv.Atoi(value)
	case MetaTypeFloat:
		parsed, err = strconv.ParseFloat(value, 64)
	case MetaTypeDateTime:
		parsed, err = time.Parse(time.RFC3339, value)
	default:
		parsed, err = value, nil
	}
	if err != nil {
		return value
	}
	return parsed
}

func MakeAddMetaHTMXButton(icon, tooltip, metaType, imageId, collection string) Node {
	url := rt.AddQueryParams(MetaUrl,
		MetaTypeQueryArg, metaType,
		MetaCollectionArg, collection,
		MetaImageIdArg, imageId,
	)
	return cmp.MakeIconizedButton(icon, tooltip,
		Attr(fmt.Sprintf("hx-get=%v", url.String())),
		Attr(fmt.Sprintf("hx-target=#%v", MetaFormId)))
}

func (s *Server) MetaTableRow(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(MetaKeyArg)
	imageId := r.URL.Query().Get(MetaImageIdArg)
	collection := r.URL.Query().Get(MetaCollectionArg)
	url := b.NewRowURL(MetaRowUrl)
	url.Set(MetaImageIdArg, imageId)
	url.Set(MetaCollectionArg, collection)
	url.Set(MetaKeyArg, key)
	switch r.URL.Query().Get("mode") {
	case b.ModeConfirmDelete.String():
		b.RenderConfirmDeleteRow(len(MetaTableFields),
			key, "key", url.Url, w)
	default:
		p := NewMetaRowPresenter(w)
		s.ReadMetaData.Execute(r.Context(),
			readmd.Request{ImageId: imageId, Collection: collection, Key: key}, p)
	}
}

func (s *Server) AddMetaData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	imageId := r.URL.Query().Get(MetaImageIdArg)
	collection := r.URL.Query().Get(MetaCollectionArg)
	valueType := r.URL.Query().Get(MetaTypeQueryArg)
	valueStr := r.FormValue(MetaValueArg)
	value := parseFormValue(valueType, valueStr)

	p := NewAddMetaPresenter(w, b.NewRowURL(MetaRowUrl))
	s.Annotator.AddMetaData.Execute(r.Context(),
		addmd.Request{
			ImageId:    imageId,
			Collection: collection,
			Key:        r.FormValue(MetaKeyArg),
			Value:      value,
		},
		&p)
	s.Annotator.ListMetaData.Execute(r.Context(),
		listmd.Request{ImageId: imageId, Collection: collection}, &p)
}

func (s *Server) MetaDataForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	url := rt.AddQueryParams(MetaUrl,
		MetaCollectionArg, r.FormValue(MetaCollectionArg),
		MetaImageIdArg, r.FormValue(MetaImageIdArg),
	)

	b := bf.NewHTMXCreateFormBuilder(url.String(), MetaFormId)
	b.AddButtonAttr(fmt.Sprintf("hx-target=#%v", MetaDivId))
	b.AddTextField(MetaKeyArg, "Key", bf.WithRequired())
	switch r.FormValue(MetaTypeQueryArg) {
	case MetaTypeString:
		b.AddSubmitQueryParam(MetaTypeQueryArg, MetaTypeString)
		b.AddTextField(MetaValueArg, "Value", bf.WithRequired())
	case MetaTypeInt:
		b.AddSubmitQueryParam(MetaTypeQueryArg, MetaTypeInt)
		b.AddIntField(MetaValueArg, "Value", bf.WithRequired())
	case MetaTypeFloat:
		b.AddSubmitQueryParam(MetaTypeQueryArg, MetaTypeFloat)
		b.AddFloatField(MetaValueArg, "Value", bf.WithRequired())
	case MetaTypeDateTime:
		b.AddSubmitQueryParam(MetaTypeQueryArg, MetaTypeDateTime)
		b.AddDateTimeField(MetaValueArg, "Value", bf.WithRequired())
	case MetaTypeBoolean:
		b.AddSubmitQueryParam(MetaTypeQueryArg, MetaTypeBoolean)
		b.AddBooleanField(MetaValueArg, "Value")
	default:
		b.AddTextField(MetaValueArg, "Value", bf.WithRequired())
	}

	Div(Class("flex items-end"), b.Build()).Render(w)
}

func (s *Server) DeleteMetaData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(MetaKeyArg)
	imageId := r.URL.Query().Get(MetaImageIdArg)
	collection := r.URL.Query().Get(MetaCollectionArg)
	s.Annotator.DeleteMetaData.Execute(r.Context(),
		delmd.Request{ImageId: imageId, Collection: collection, Key: key},
		NewMetaDeletePresenter(w, b.RowURL{}),
	)
}
