package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lejeunel/go-image-annotator/config"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
	"github.com/spf13/cobra"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

type DoctorCLIPresenter struct {
	NumErrors int
	NumExams  int
}

func (p *DoctorCLIPresenter) PresentExam(r dr.Response) {
	p.NumExams += 1

	if r.Error == nil {
		fmt.Printf("%s: %s\n", r.Description, green("OK"))
		return
	}
	p.NumErrors += 1
	fmt.Printf("%s: %s %v\n", r.Description, red("ERROR"), r.Error)
}

// Summarize prints a summary of exams and error counts.
func (p *DoctorCLIPresenter) Summarize() {
	fmt.Printf("Performed %v exams and got %v error(s)\n", p.NumExams, p.NumErrors)
}

type FromEnvConfigProvider struct{}

func (p FromEnvConfigProvider) Get() (*config.Config, error) {
	return config.Parse()
}

type PasswordValidatorProvider struct{}

func (p PasswordValidatorProvider) Get(minEntropy int) pv.PasswordValidator {
	return pv.New(minEntropy)
}

type FileStoreProvider struct{}

func (p FileStoreProvider) Get(c config.Config) (fs.FileStore, error) {
	return fs.Build(c)
}

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks configuration variables and output a diagnoses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		itr := dr.New(FromEnvConfigProvider{}, PasswordValidatorProvider{},
			FileStoreProvider{})
		p := &DoctorCLIPresenter{}
		itr.Execute(context.Background(), p)
		p.Summarize()
		if p.NumErrors > 0 {
			os.Exit(1)
		}
	},
}
