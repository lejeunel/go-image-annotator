package annotator

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
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

type AddMetaPresenter struct {
	data   []m.MetaData
	Writer http.ResponseWriter
	GetURL b.RowURL
	htmx.ErrorPresenter
}

func NewAddMetaPresenter(w http.ResponseWriter, u b.RowURL) AddMetaPresenter {
	return AddMetaPresenter{
		Writer:         w,
		GetURL:         u,
		ErrorPresenter: htmx.NewErrorPresenter("adding meta-data", w),
	}
}

func (p *AddMetaPresenter) SuccessAddMetadata() {}
func (p *AddMetaPresenter) SuccessListMetadata(r listmd.Response) {
	RenderMetaDataList(r.ImageId, r.Collection, r.MetaData, p.Writer)
}

type MetaPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	htmx.ErrorPresenter
}

func NewMetaPresenter(w http.ResponseWriter, u b.RowURL) MetaPresenter {
	return MetaPresenter{w, u, htmx.NewErrorPresenter("viewing meta-data", w)}
}

func (p MetaPresenter) SuccessReadMetadata(r readmd.Response) {
	MakeRow(r.ImageId, r.Collection, r.Data).Render(p.writer)
}
func (p MetaPresenter) SuccessDeleteMetadata(key string) {
	htmx.NotifySuccessPayload(p.writer,
		"deleting meta-data item",
		fmt.Sprintf("Successfully deleted key %v", key))
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
		p := NewMetaPresenter(w, url)
		s.ReadMetaData.Execute(r.Context(),
			readmd.Request{ImageId: imageId, Collection: collection, Key: key}, p)
	}
}

func (s *Server) AddMetaData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	p := NewAddMetaPresenter(w, b.NewRowURL(MetaRowUrl))
	imageId := r.URL.Query().Get(MetaImageIdArg)
	collection := r.URL.Query().Get(MetaCollectionArg)
	valueType := r.URL.Query().Get(MetaTypeQueryArg)
	valueStr := r.FormValue(MetaValueArg)
	value := parseFormValue(valueType, valueStr)

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
		NewMetaPresenter(w, b.RowURL{}),
	)
}

func MakeRow(imageId, collection string, data m.MetaData) tb.Row {
	url := b.NewRowURL(MetaRowUrl)
	url.Set(MetaImageIdArg, imageId)
	url.Set(MetaCollectionArg, collection)
	url.Set(MetaKeyArg, data.Key)

	actions := b.NewActionsPanelBuilder()
	actions.SetConfirmDelete(url.SetMode(b.ModeConfirmDelete).Url)
	r := tb.NewRow()
	r.AddCell(tb.NewCell(Text(data.Key)))
	r.AddCell(tb.NewCell(Text(fmt.Sprint(data.Value))))
	r.AddCell(tb.NewCell(actions.Build()))
	return r
}

func BuildMetaDataList(imageId, collection string, data []m.MetaData) Node {
	table := tb.NewTableBuilder(MetaTableFields, tb.WithSimplePlaceHolder())
	for _, d := range data {
		table.AddRow(MakeRow(imageId, collection, d))
	}

	buttons := Div(
		Class("flex gap-1"),
		Group([]Node{
			MakeAddMetaHTMXButton(
				ic.Text,
				"Add text",
				MetaTypeString,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Hash,
				"Add integer",
				MetaTypeInt,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Percent,
				"Add float",
				MetaTypeFloat,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Calendar,
				"Add date/time",
				MetaTypeDateTime,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Flag,
				"Add boolean",
				MetaTypeBoolean,
				imageId,
				collection,
			),
		}))
	return Div(
		ID(MetaDivId),
		cmp.Separator,
		Div(Class("flex items-center mt-2 mb-2"),
			Div(Class("flex-1 text-lg"), Text("Custom meta-data")),
			buttons,
		),
		Div(ID(MetaFormId)),
		table.Build())
}

func RenderMetaDataList(imageId, collection string, data []m.MetaData, w io.Writer) {
	BuildMetaDataList(imageId, collection, data).Render(w)
}
