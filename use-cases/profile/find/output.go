package find

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type OutputPort interface {
	Error(error)
	SuccessFindProfile(pr.Profile)
}
