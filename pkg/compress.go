package pkg

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/pkg/compressor"
	"github.com/ONSDigital/blaise-mi-extract/pkg/storage/google"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var compressDestination string
var zipOnce sync.Once

func initialiseZip() {
	util.Initialise()

	var found bool

	if compressDestination, found = os.LookupEnv(util.ZipLocation); !found {
		log.Fatal().Msg("The " + util.ZipLocation + " variable has not been set")
		os.Exit(1)
	}

	log.Info().
		Str("location", compressDestination).Msg("Zip Destination")
}

// handles event from item arriving in the compress bucket
func ZipCompress(ctx context.Context, name, location string) error {

	zipOnce.Do(func() {
		initialiseZip()
	})

	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received compress request")

	r := google.NewStorage(ctx)
	zip := compressor.NewService(&r)

	var zipName string
	var err error

	if zipName, err = zip.ZipFile(name, location, compressDestination); err != nil {
		log.Err(err).Msg("create compress failed")
		return err
	}

	if err := zip.DeleteFile(name, location); err != nil {
		return err
	}

	log.Info().Msgf("file %s zipped and saved to %s/%s", name, compressDestination, zipName)

	return nil
}
