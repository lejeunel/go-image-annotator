package dashboard

import (
	_ "embed"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	"github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	rt "github.com/lejeunel/go-image-annotator/routes"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/log/list"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed logs-preamble.md
var logsPreamble string

var listEventsFields = []string{"id", "type", "started", "status", "actions"}

func (s *Server) ListTasks(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListTasksItr.Execute(r.Context(),
		pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewTaskListPresenter(w, s.PageBuilder))
}

type TaskListPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.ErrorPresenter
}

func NewTaskListPresenter(w http.ResponseWriter, p b.PageBuilder) TaskListPresenter {
	p.SetTitle(LogsPageName)
	p.SetHTMLTitle(LogsPageName)
	p.SetActiveSection(cmp.NoPageActive)
	p.ActivateSidebarEntry(LogsPageName)
	p.AddMarkdownPreamble(logsPreamble)
	b := b.NewPaginatedListBuilder(p, listEventsFields)
	return TaskListPresenter{b, w, e.NewErrorPresenter(w)}
}

func (p TaskListPresenter) SuccessListTasks(r list.Response) {
	for _, t := range r.Tasks {
		row := MakeRow(t)
		p.AddRow(row)
	}
	p.SetPagination(r.Pagination, rt.ListTasksUrl)
	p.Render(p.Writer)
}

type TaskRowPresenter struct {
	Writer io.Writer
	e.ErrorPresenter
}

func NewTaskRowPresenter(w http.ResponseWriter) TaskRowPresenter {
	return TaskRowPresenter{w, e.NewErrorPresenter(w)}
}

func (p TaskRowPresenter) SuccessFindTask(t t.Task) {
	MakeRow(t).Render(p.Writer)
}

func MakeRow(t t.Task) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetExpand(rt.AddQueryParams(TaskDetailsUrl, TaskIdQueryArg,
		t.Id.String()))

	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(t.Id.String())))
	row.AddCell(tb.NewCell(Text(t.Type.String())))
	if len(t.Events) > 0 {
		row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(t.Events[len(t.Events)-1].Time))))
	}

	var state Node
	switch t.Events[0].State {
	case event.PendingTask:
		state = Div(Class("text-warning"), Text(event.PendingTask.String()))
	case event.StartedTask:
		state = Div(Class("text-warning"), Text(event.StartedTask.String()))
	case event.FailedTask:
		state = Div(Class("text-danger"), Text(event.FailedTask.String()))
	default:
		state = Div(Class("text-success"), Text(event.DoneTask.String()))
	}
	row.AddCell(tb.NewCell(state))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
