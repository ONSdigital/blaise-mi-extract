package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/zipper"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/google"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
)

var encryptedDestination string
var zip zipper.Service

func init() {
	util.Initialise()

	r := google.NewStorage()
	zip = zipper.NewService(r)

	var found bool

	if encryptedDestination, found = os.LookupEnv(util.EncryptedLocation); !found {
		log.Fatal().Msg("The " + util.EncryptedLocation + " variable has not been set")
		os.Exit(1)
	}
}

func ZipFunction(_ context.Context, e util.GCSEvent) error {

	log.Info().
		Str("bucket", e.Bucket).
		Str("file", e.Name).
		Msgf("received zip request")

	var zipName string
	var err error

	if zipName, err = zip.ZipFile(e.Name, e.Bucket, encryptedDestination); err != nil {
		log.Err(err).Msg("create zip failed")
		return err
	}

	if err := zip.DeleteFile(e.Name, e.Bucket); err != nil {
		return err
	}

	log.Info().Msgf("file %s zipped and saved to %s/%s", e.Name, encryptedDestination, zipName)

	return nil
}
