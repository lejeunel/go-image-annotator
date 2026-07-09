package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var createCollectionTargetDiv = "create-collection"

type CreateCollectionPresenter struct {
	w http.ResponseWriter
}

func (p CreateCollectionPresenter) Success(r create.Response) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify-and-reload": map[string]string{
			"variant": "success",
			"title":   "created collection",
			"message": fmt.Sprintf("Successfully create collection with name %v", r.Name),
		},
	})
	p.w.Header().Set("HX-Trigger", string(payload))
	p.w.WriteHeader(http.StatusOK)
}
func (p CreateCollectionPresenter) Error(err error) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify": map[string]string{
			"variant": "danger",
			"title":   "failed creating collection",
			"message": err.Error(),
		},
	})
	p.w.Header().Set("HX-Trigger", string(payload))
	p.w.WriteHeader(http.StatusUnprocessableEntity)
}

type ListCollectionsPresenter struct {
	ListRenderer
}

func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	table := html.MyTable{Fields: []string{"name", "description", "group", "created", "actions"}}
	for _, c := range r.Collections {
		var groupName string
		if c.Group == nil {
			groupName = "n/a"
		} else {
			groupName = c.Group.Name
		}

		actions := html.NewActionsPanel()
		actions.SetEdit("/edit-url")
		actions.SetDelete("/delete-url")
		table.AddRow(html.MyTableRow{Values: []Node{html.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name),
			Raw(c.Description), Raw(groupName), Raw(cmp.DateTimeToStr(c.CreatedAt)), actions.Build()}})
	}
	button := cmp.MakeHTMXCreateButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	preamble := Div(ID(createCollectionTargetDiv))
	p.RenderList(&preamble, table, r.Pagination, &button)
}
func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.Collection.Create.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, CreateCollectionPresenter{w})

}

func (s *Server) CreateCollectionForm(w http.ResponseWriter, r *http.Request) {

	form := Span(Class("w-full inline-flex items-center justify-end"),
		MakeCreateCollectionForm(createCollectionTargetDiv))

	form.Render(w)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}

func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	baseURL, _ := url.Parse("/collections")
	return ListCollectionsPresenter{
		ListRenderer: NewListRenderer(*p.SetTitle("Collections").SetActive(b.CollectionsPageActive), *baseURL,
			w),
	}
}

func MakeCreateCollectionForm(containerId string) Node {
	return Form(
		Attr(fmt.Sprintf(`hx-post=%v`, rt.CreateCollection)),
		Attr(`hx-swap="none"`),
		Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
		Label(For("name"), Text("Name"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
		Input(Type("text"), ID("email"), Name("name"), Required(),
			Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
		Label(For("description"), Text("Description"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
		Input(Type("text"), ID("description"), Name("description"),
			Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
		Span(Class("flex items-center"),
			Button(Type("submit"),
				Text("Submit"),
				Class(st.ActivePrimaryButton)),
			Button(Type("button"),
				Text("Cancel"),
				Class(st.InactiveButton),
				Attr(`hx-on:click`, fmt.Sprintf(`document.getElementById('%v').innerHTML=''`, containerId))),
		),
	)
}
