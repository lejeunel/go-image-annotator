package dashboard

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type TaskDetailPresenter struct {
	io.Writer
	htmx.ErrorPresenter
}

func NewTaskDetailPresenter(w http.ResponseWriter) TaskDetailPresenter {
	task := "detailing task"
	return TaskDetailPresenter{w, htmx.NewErrorPresenter(task, w)}
}

type Task struct {
	Id     string  `json:"id"`
	Type   string  `json:"type"`
	Issuer string  `json:"issuer"`
	Events []Event `json:"events"`
}
type Event struct {
	Time  string            `json:"time"`
	State string            `json:"state"`
	Extra map[string]string `json:"extra"`
	Error string            `json:"error"`
}

func (p *TaskDetailPresenter) SuccessFindTask(t t.Task) {
	r := Task{Id: t.Id.String(), Type: t.Type.String(), Issuer: t.Issuer}
	for _, e := range t.Events {
		r.Events = append(r.Events,
			Event{Time: e.Time.Format(time.DateTime),
				State: e.State.String(), Extra: e.Extra, Error: e.Error,
			})
	}
	data, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		Text(err.Error()).Render(p.Writer)
		return
	}
	url := rt.AddQueryParams(TaskRowUrl, TaskIdQueryArg, t.Id.String())
	Tr(Td(Attr("colspan=3"), Class("p-4"),
		Div(Pre(Class("bg-surface-alt dark:bg-surface-dark-alt p-2"),
			Text(string(data))))),
		Td(Class("align-top p-4"),
			cmp.MakeHTMXAbortButton("Close", url.String()),
		)).Render(p.Writer)
}

func (s *Server) TaskRow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(TaskIdQueryArg)
	p := NewTaskRowPresenter(w)
	s.FindTaskItr.Execute(r.Context(), id, &p)
}

func (s *Server) TaskDetails(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(TaskIdQueryArg)
	p := NewTaskDetailPresenter(w)
	s.FindTaskItr.Execute(r.Context(), id, &p)
}
