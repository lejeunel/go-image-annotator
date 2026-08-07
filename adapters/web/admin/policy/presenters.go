package policy

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	_ "embed"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	a "github.com/lejeunel/go-image-annotator/modules/authorizer"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed preamble.md
var preambleTemplate string

//go:embed postamble.md
var postambleTemplate string

type SetPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(string) string
	htmx.ErrorPresenter
}

func NewSetPresenter(w http.ResponseWriter) SetPresenter {
	task := "Updating access policies"
	okMessageFunc := func(string) string {
		return "Successfully updated policies"
	}
	return SetPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p SetPresenter) SuccessSetPolicy(policies string) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(policies))
}

type PostAmbleData struct {
	Methods []string
}
type PreAmbleData struct {
	DefaultPoliciesURL string
}

type ViewPresenter struct {
	b.PageBuilder
	io.Writer
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, p b.PageBuilder) ViewPresenter {
	return ViewPresenter{p, w, e.NewErrorPresenter(w)}
}

func (p ViewPresenter) SuccessReadPolicy(policies string) {
	t, err := template.New("").Parse(postambleTemplate)
	if err != nil {
		Text(err.Error()).Render(p.Writer)
		return
	}
	var bufPost bytes.Buffer
	if err := t.Execute(&bufPost, PostAmbleData{Methods: a.ValidMethods}); err != nil {
		Text(err.Error()).Render(p.Writer)
		return
	}

	t, err = template.New("").Parse(preambleTemplate)
	if err != nil {
		Text(err.Error()).Render(p.Writer)
		return
	}
	var bufPre bytes.Buffer
	if err := t.Execute(&bufPre, PreAmbleData{DefaultPolicyDownloadUrl}); err != nil {
		Text(err.Error()).Render(p.Writer)
		return
	}

	p.AddMarkdownPostamble(bufPost.String())
	p.AddMarkdownPreamble(bufPre.String())
	textArea := Textarea(
		Name(PolicyFieldName),
		Class(`w-150 h-90 rounded-lg border-2 border-surface-alt dark:border-surface-dark-alt p-3
         focus:outline-none focus:border-primary focus:dark:border-primary-dark
         resize-none`),
		Text(policies))
	form := Form(
		Attr(fmt.Sprintf(`hx-post=%v`, SetPolicyFormUrl)),
		textArea,
		Span(Class("flex mt-4"),
			Button(Type("submit"),
				Text("Submit"),
				Class(st.SuccessButton)),
		),
	)
	p.SetContent(form)
	p.Render(p.Writer)
}
