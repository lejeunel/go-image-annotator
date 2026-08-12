package image

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	s "github.com/lejeunel/go-image-annotator/adapters/web/shared"
	rt "github.com/lejeunel/go-image-annotator/routes"
	ia "github.com/lejeunel/go-image-annotator/use-cases/image/ingest-archive"
)

//go:embed templates/ingest.html
var IngestionPanel string

//go:embed templates/ingest.py
var PythonIngestionScript string

type IngestArchivePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(ia.Response) string
	htmx.ErrorPresenter
}

func NewIngestArchivePresenter(w http.ResponseWriter) IngestArchivePresenter {
	task := "Ingesting image archive"
	okMessageFunc := func(r ia.Response) string {
		return s.MakeNewTaskMessage()
	}
	return IngestArchivePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p IngestArchivePresenter) SuccessSubmitIngestArchiveTask(r ia.Response) {
	htmx.NotifySuccessPayload(p.writer, p.task, p.okMessageFunc(r))
}

type IngestionPanelData struct {
	DivId                 string
	ArchiveIngestUrl      string
	PythonIngestionScript string
	InputName             string
	MaxMB                 int
}

func (s *Server) IngestionPanel(w http.ResponseWriter, r *http.Request) {
	t := template.New("")
	t.Parse(IngestionPanel)

	endpoint := rt.AddQueryParams(archiveIngestUrl,
		ingestCollectionArgName, r.FormValue(ingestCollectionArgName))
	if err := t.ExecuteTemplate(w, "ingestion",
		IngestionPanelData{
			DivId:                 ingestTargetDiv,
			ArchiveIngestUrl:      endpoint.String(),
			MaxMB:                 s.maxArchiveMB,
			InputName:             ingestFormInputName,
			PythonIngestionScript: PythonIngestionScript,
		}); err != nil {
		panic(err)
	}
}

func (s *Server) IngestArchive(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(ingestFormInputName)
	if err != nil {
		http.Error(w, "missing archive: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	s.IngestArchiveItr.Execute(r.Context(),
		ia.Request{Collection: r.URL.Query().Get(ingestCollectionArgName),
			Reader: file,
		},
		NewIngestArchivePresenter(w),
	)
}
