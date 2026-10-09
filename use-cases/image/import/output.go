package import_image

type OutputPort interface {
	Error(error)
	SuccessImportImage(Response)
}
