package cmd

import (
	"context"
	"os"

	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
	"github.com/spf13/cobra"
)

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks configuration variables and output a diagnoses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		itr := dr.New(FromEnvConfigProvider{}, PasswordValidatorProvider{},
			FileStoreProvider{}, DBStoreProvider{}, ImageStoreProvider{})
		p := NewDoctorCLIPresenter()
		itr.Execute(context.Background(), p)
		p.Summarize()
		if p.NumErrors > 0 {
			os.Exit(1)
		}
	},
}
