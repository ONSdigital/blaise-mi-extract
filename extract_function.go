package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/extractor"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/google"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
)

type PubSubMessage struct {
	Action     string `json:"action"`
	Instrument string `json:"instrument_name"`
}

var encryptDestination string
var service extractor.Service

func init() {
	util.Initialise()

	r := google.NewStorage()
	service = extractor.NewService(r)

	var found bool

	if encryptDestination, found = os.LookupEnv(encryptLocation); !found {
		log.Fatal().Msg("The " + encryptLocation + " varible has not been set")
		os.Exit(1)
	}
}

func ExtractFunction(_ context.Context, m PubSubMessage) error {

	if m.Action != "extract_mi" {
		log.Warn().Msgf("message rejected, unknown action -> [%s]", m.Action)
		return nil
	}

	log.Info().
		Str("action", m.Action).
		Str("instrument", m.Instrument).
		Msgf("received extract request")

	var source string
	var err error

	if source, err = service.DataToCSV(); err != nil {
		return err
	}

	destination := m.Instrument + ".csv"

	if err = service.SaveFile(encryptDestination, source, destination); err != nil {
		return err
	}

	return nil
}
