package create

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type OutputPort interface {
	SuccessCreateProfile(pr.Profile)
	Error(error)
}
