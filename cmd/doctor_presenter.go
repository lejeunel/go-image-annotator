package cmd

import (
	"fmt"

	"github.com/fatih/color"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
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

func (p *DoctorCLIPresenter) SuccessCheckImageIntegrity(id im.ImageId) {
	fmt.Println("success check image integrity:", id)
}

func (p *DoctorCLIPresenter) ErrorCheckImageIntegrity(id im.ImageId) {
	fmt.Println("error check image integrity:", id)
}

// Summarize prints a summary of exams and error counts.
func (p *DoctorCLIPresenter) Summarize() {
	fmt.Printf("Performed %v exams and got %v error(s)\n", p.NumExams, p.NumErrors)
}
