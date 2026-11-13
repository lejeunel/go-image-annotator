package doctor

import im "github.com/lejeunel/go-image-annotator/entities/image"

type OutputPort interface {
	PresentExam(Response)
	SuccessCheckImageIntegrity(im.ImageId)
	ErrorCheckImageIntegrity(im.ImageId)
	// DoneCheckImageIntegrity is called once the integrity check is over,
	// before the error that interrupted it, if any, is presented.
	DoneCheckImageIntegrity()
}
