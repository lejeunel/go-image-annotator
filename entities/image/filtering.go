package image

type OrderingParams struct {
	IngestTime bool
}

type FilteringParams struct {
	Collection *string
}

type CountingParams struct {
	Collection *string
}
