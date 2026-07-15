package image

type OrderingParams struct {
	IngestTime bool
}

type FilteringParams struct {
	Collection *string
}

type PaginationParams struct {
	PageSize int
	Page     int64
}

type CountingParams struct {
	Collection *string
}
