package sqlite

import (
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
)

func NewDoctorInteractor(pv pw.PasswordValidator,
) dr.Interactor {
	return dr.New(pv)
}
