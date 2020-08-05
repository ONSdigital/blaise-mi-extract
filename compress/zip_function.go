package compress

import (
	"github.com/ONSDigital/blaise-mi-extract/pkg/zipper"
	"github.com/ONSDigital/blaise-mi-extract/storage/google"
	"github.com/ONSDigital/blaise-mi-extract/util"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var encryptedDestination string
var zip zipper.Service
var doOnce sync.Once

func initialise() {
	util.Initialise()

	r := google.NewStorage()
	zip = zipper.NewService(r)

	var found bool

	if encryptedDestination, found = os.LookupEnv(util.EncryptedLocation); !found {
		log.Fatal().Msg("The " + util.EncryptedLocation + " variable has not been set")
		os.Exit(1)
	}

	log.Info().
		Str("location", encryptedDestination).Msg("Encryption Destination")

}

// handles event from item arriving in the compress bucket
func ZipCompress(name, location string) error {

	doOnce.Do(func() {
		initialise()
	})

	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received compress request")

	var zipName string
	var err error

	if zipName, err = zip.ZipFile(name, location, encryptedDestination); err != nil {
		log.Err(err).Msg("create compress failed")
		return err
	}

	if err := zip.DeleteFile(name, location); err != nil {
		return err
	}

	log.Info().Msgf("file %s zipped and saved to %s/%s", name, encryptedDestination, zipName)

	return nil
}
