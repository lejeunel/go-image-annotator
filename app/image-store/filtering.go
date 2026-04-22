package image_store

type FilteringParams struct {
	Collection *string
	PageSize   int
	Page       int64
}

type CountingParams struct {
	Collection *string
}
