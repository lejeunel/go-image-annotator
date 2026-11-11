package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/lejeunel/go-image-annotator/config"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
	"github.com/spf13/cobra"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

type DoctorCLIPresenter struct{}

func (p DoctorCLIPresenter) PresentExam(r dr.Response) {
	if r.Error == nil {
		fmt.Printf("%s: %s\n", r.Description, green("OK"))
		return
	}
	fmt.Printf("%s: %s %v\n", r.Description, red("ERROR"), r.Error)
}

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks configuration variables and output a diagnoses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Parse()
		itr := dr.NewFromConfig(cfg)
		itr.Execute(context.Background(),
			dr.Request{Config: cfg},
			DoctorCLIPresenter{})
	},
}
