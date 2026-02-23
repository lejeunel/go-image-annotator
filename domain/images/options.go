package images

type FetchImageOptions struct {
	IncludeRawData bool
}

type ImportImageOptions struct {
	ImportAnnotations bool
}

var ImportImageWithAnnotations = ImportImageOptions{ImportAnnotations: true}
var ImportImageWithoutAnnotations = ImportImageOptions{ImportAnnotations: false}

var FetchMetaOnly = FetchImageOptions{IncludeRawData: false}
var FetchWithRawData = FetchImageOptions{IncludeRawData: true}
