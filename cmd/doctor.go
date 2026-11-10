package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/fatih/color"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
	"github.com/spf13/cobra"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

type DoctorCLIPresenter struct{}

func (p DoctorCLIPresenter) SuccessMeetTheDoctor(r dr.Response) {
	for _, name := range slices.Sorted(maps.Keys(r.Exams)) {
		result := r.Exams[name]
		if result.Error == nil {
			fmt.Printf("%s: %s\n", result.Description, green("OK"))
			continue
		}
		fmt.Printf("%s: %s %v\n", result.Description, red("ERROR"), result.Error)
		if result.Advice != "" {
			fmt.Printf("  advice: %s\n", result.Advice)
		}
	}
}

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks configuration variables and output a diagnoses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.DiscardHandler)
		app := sqlite.NewAppFromEnv(logger)
		app.Itrs.Doctor.Execute(context.Background(),
			dr.Request{InitialAdminPassword: app.Config.InitialAdminPassword},
			DoctorCLIPresenter{})
	},
}
