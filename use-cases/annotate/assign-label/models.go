package assign_label

type Response struct {
	AnnotationId string
	ImageId      string
	Collection   string
	Label        string
}

type Request struct {
	ImageId    string
	Collection string
	Label      string
}
