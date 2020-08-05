package encrypt

import (
	"github.com/ONSDigital/blaise-mi-extract/pkg/encryption"
	"github.com/ONSDigital/blaise-mi-extract/storage/google"
	"github.com/ONSDigital/blaise-mi-extract/util"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var zipDestination string
var keyFile string
var doOnce sync.Once
var encrypter encryption.Service

func initialise() {
	util.Initialise()

	r := google.NewStorage()
	encrypter = encryption.NewService(r)

	var found bool

	if zipDestination, found = os.LookupEnv(util.ZipLocation); !found {
		log.Fatal().Msg("The " + util.ZipLocation + " variable has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("compress destination: %s", zipDestination)

	if keyFile, found = os.LookupEnv(util.PublicKeyFile); !found {
		log.Fatal().Msg("The " + util.PublicKeyFile + " variable has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("public key file: %s", keyFile)
}

// handles event from item arriving in the encrypt  bucket
func HandleEncryptionRequest(name, location string) error {

	doOnce.Do(func() {
		initialise()
	})

	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received encrypt request")

	if err := encrypter.EncryptFile(keyFile, name, location, zipDestination); err != nil {
		log.Err(err).Msg("encrypt failed")
		return err
	}

	if err := encrypter.DeleteFile(name, location); err != nil {
		return err
	}

	return nil
}
