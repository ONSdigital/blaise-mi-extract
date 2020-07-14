package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/encryption"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/google"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
)

var zipDestination string
var keyFile string

var encrypter encryption.Service

func init() {
	util.Initialise()

	r := google.NewStorage()
	encrypter = encryption.NewService(r)

	var found bool

	if zipDestination, found = os.LookupEnv(zipLocation); !found {
		log.Fatal().Msg("The " + zipLocation + " varible has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("zip destination: %s", zipDestination)

	if keyFile, found = os.LookupEnv(publicKeyFile); !found {
		log.Fatal().Msg("The " + publicKeyFile + " varible has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("public key file: %s", keyFile)
}

func EncryptFunction(_ context.Context, e GCSEvent) error {

	log.Info().
		Str("bucket", e.Bucket).
		Str("file", e.Name).
		Msgf("received encrypt request")

	if err := encrypter.EncryptFile(keyFile, e.Name, e.Bucket, zipDestination); err != nil {
		log.Err(err).Msg("encrypt failed")
		return err
	}

	if err := encrypter.DeleteFile(e.Name, e.Bucket); err != nil {
		return err
	}

	return nil
}
