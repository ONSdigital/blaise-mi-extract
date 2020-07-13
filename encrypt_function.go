package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/encryption"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/google"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
)

var encryptedDestination string
var keyFile string

var encrypter encryption.Service

func init() {
	util.Initialise()

	r := google.NewStorage()
	encrypter = encryption.NewService(r)

	var found bool

	if encryptedDestination, found = os.LookupEnv(encryptedLocation); !found {
		log.Fatal().Msg("The " + encryptedLocation + " varible has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("encrypted destination: %s", encryptedDestination)

	if keyFile, found = os.LookupEnv(publicKeyFile); !found {
		log.Fatal().Msg("The " + publicKeyFile + " varible has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("public key file: %s", keyFile)
}

func EncryptFunction(ctx context.Context, e GCSEvent) error {

	log.Info().
		Str("bucket", e.Bucket).
		Str("file", e.Name).
		Msgf("received encrypt request")

	if err := encrypter.EncryptFile(keyFile, e.Name, e.Bucket, encryptedDestination); err != nil {
		log.Err(err).Msg("encrypt failed")
		return err
	}

	if err := encrypter.DeleteFile(e.Name, e.Bucket); err != nil {
		return err
	}

	log.Info().Msgf("file %s encrypted and saved to %s", e.Name, encryptedDestination)

	return nil
}
