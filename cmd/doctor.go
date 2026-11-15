package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	ev "github.com/lejeunel/go-image-annotator/modules/email-validator"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
	"github.com/spf13/cobra"
)

var (
	missingImagesLog string
	DoctorCmd        = &cobra.Command{
		Use:   "doctor",
		Short: "Checks configuration variables and output a diagnoses",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			itr := dr.New(FromEnvConfigProvider{}, PasswordValidatorProvider{}, ev.New(),
				FileStoreProvider{}, DBStoreProvider{}, ImageStoreProvider{}, NotifierProvider{})
			p := NewDoctorCLIPresenter(missingImagesLog)
			itr.Execute(context.Background(), p)
			p.Summarize()
			if p.NumErrors > 0 {
				os.Exit(1)
			}
		},
	}
)

func init() {
	defaultLog := fmt.Sprintf(".doctor-missing-images-%s.log",
		time.Now().Format("20060102T150405"))
	DoctorCmd.Flags().StringVar(&missingImagesLog, "missing-images-log", defaultLog,
		"file name where missing image IDs are written")
}
