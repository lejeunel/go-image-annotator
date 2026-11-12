package doctor

import im "github.com/lejeunel/go-image-annotator/entities/image"

type OutputPort interface {
	PresentExam(Response)
	SuccessCheckImageIntegrity(im.ImageId)
	ErrorCheckImageIntegrity(im.ImageId)
}
