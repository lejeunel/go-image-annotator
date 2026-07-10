package web

import (
	"fmt"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
)

var createCollectionTargetDiv = "create-collection"

type CreateCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	WebErrorPresenter
}

func NewCreateCollectionPresenter(w http.ResponseWriter) CreateCollectionPresenter {
	task := "Creating collection"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created collection with name %v", r.Name)
	}
	return CreateCollectionPresenter{w, task, okMessageFunc, NewWebErrorPresenter(task, w)}
}
func (p CreateCollectionPresenter) Success(r create.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Collection.Create.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, NewCreateCollectionPresenter(w))
}
func (s *Server) CreateCollectionForm(w http.ResponseWriter, r *http.Request) {

	form := Span(Class("w-full inline-flex items-center justify-end"),
		MakeCreateCollectionForm(createCollectionTargetDiv))

	form.Render(w)
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
