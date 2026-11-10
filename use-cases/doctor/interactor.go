package doctor

import (
	"context"

	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
)

type Interactor struct {
	pw.PasswordValidator
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	response := NewDoctorResponse()

	adminpw := ExamResult{
		Description: "Initial administrator password",
		Error:       i.PasswordValidator.Validate(r.InitialAdminPassword),
	}
	response.AddExam(InitialAdminPassword, adminpw)

	out.SuccessMeetTheDoctor(response)
}

type Option func(*Interactor)

func New(passwordValidator pw.PasswordValidator) Interactor {
	return Interactor{
		PasswordValidator: passwordValidator,
	}
}
