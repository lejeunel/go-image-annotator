package ingest

type OutputPort interface {
	SuccessSubmitIngestArchiveTask(Response)
	Error(error)
}
