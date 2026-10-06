package update

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type OutputPort interface {
	SuccessUpdateProfile(pr.Profile)
	Error(error)
}
