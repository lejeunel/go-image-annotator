package dashboard

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/log/list"
	. "maragu.dev/gomponents"
)

//go:embed logs-preamble.md
var logsPreamble string

var listEventsFields = []string{"task_id", "task_type", "time", "state", "extra", "error"}

func (s *Server) ListTasks(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListTasksItr.Execute(r.Context(),
		pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewLogsPresenter(w, s.PageBuilder))
}

type LogsPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.ErrorPresenter
}

func NewLogsPresenter(w http.ResponseWriter, p b.PageBuilder) LogsPresenter {
	p.SetTitle(LogsPageName)
	p.SetHTMLTitle(LogsPageName)
	p.SetActiveSection(cmp.NoPageActive)
	p.ActivateSidebarEntry(LogsPageName)
	p.AddMarkdownPreamble(logsPreamble)
	b := b.NewPaginatedListBuilder(p, listEventsFields)
	return LogsPresenter{b, w, e.NewErrorPresenter(w)}
}
func (p LogsPresenter) SuccessListTasks(r list.Response) {
	for _, t := range r.Tasks {
		rows := MakeRows(t)
		for _, r := range rows {
			p.AddRow(r)
		}
	}
	fmt.Printf("%+v\n", r.Pagination)
	p.SetPagination(r.Pagination, LogsUrl)
	p.Render(p.Writer)
}

func mapToJSON(m map[string]string) (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func MakeRows(t t.Task) []tb.Row {
	var rows []tb.Row
	for _, ev := range t.Events {
		row := tb.NewRow()
		row.AddCell(tb.NewCell(Text(t.Id.String())))
		row.AddCell(tb.NewCell(Text(t.Type.String())))
		row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(ev.Time))))
		row.AddCell(tb.NewCell(Text(ev.State.String())))
		extra, err := mapToJSON(ev.Extra)
		if err != nil {
			row.AddCell(tb.NewCell(Text(err.Error())))
		} else {
			row.AddCell(tb.NewCell(Text(extra)))
		}
		row.AddCell(tb.NewCell(Text(ev.Error)))
		rows = append(rows, row)

	}
	return rows

}
