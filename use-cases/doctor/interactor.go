package doctor

import (
	"context"
	"fmt"

	"github.com/lejeunel/go-image-annotator/config"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
)

type ConfigProvider interface {
	Get() (*config.Config, error)
}

type PasswordValidatorProvider interface {
	Get(minEntropy int) pv.PasswordValidator
}

type Interactor struct {
	PasswordValidatorProvider
	ConfigProvider
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	config, err := i.ConfigProvider.Get()
	parseResponse := Response{Description: "Parsing configuration"}
	if err != nil {
		parseResponse.Error = err
		out.PresentExam(parseResponse)
		return
	}
	out.PresentExam(parseResponse)

	pwValidator := i.PasswordValidatorProvider.Get(config.PasswordMinEntropy)
	out.PresentExam(Response{
		Description: "Initial administrator password",
		Error:       pwValidator.Validate(config.InitialAdminPassword),
	})
	fsType, _, err := fs.ParseURI(config.ImagesURI)
	out.PresentExam(Response{
		Description: fmt.Sprintf("Image store URI (%v)", fsType),
		Error:       err,
	})

	if fsType == fs.S3FileStoreType {
		_, err := config.ParseS3Config(ctx)
		out.PresentExam(Response{
			Description: "Validating S3 bucket parameters",
			Error:       err,
		})
	}
}

type Option func(*Interactor)

func New(cfgp ConfigProvider, pvp PasswordValidatorProvider) Interactor {
	return Interactor{
		pvp, cfgp,
	}
}
