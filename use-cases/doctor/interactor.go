package doctor

import (
	"context"
	"fmt"

	"github.com/lejeunel/go-image-annotator/config"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
)

type Interactor struct {
	pw.PasswordValidator
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	out.PresentExam(Response{
		Name:        InitialAdminPassword,
		Description: "Initial administrator password",
		Error:       i.PasswordValidator.Validate(r.Config.InitialAdminPassword),
	})

	fsType, _, err := fs.ParseURI(r.Config.ImagesURI)
	out.PresentExam(Response{
		Name:        ImageStoreURI,
		Description: fmt.Sprintf("Image store URI (%v)", fsType),
		Error:       err,
	})

	if fsType == fs.S3FileStoreType {
		_, err := r.Config.ParseS3Config(ctx)
		out.PresentExam(Response{
			Name:        S3ImageStoreConfig,
			Description: "Validating S3 bucket parameters",
			Error:       err,
		})
	}
}

type Option func(*Interactor)

func New(passwordValidator pw.PasswordValidator) Interactor {
	return Interactor{
		PasswordValidator: passwordValidator,
	}
}

func NewFromConfig(cfg config.Config) Interactor {
	passwordValidator := pv.New(cfg.PasswordMinEntropy)
	return New(passwordValidator)
}
