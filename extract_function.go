package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/extractor"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/google"
	"github.com/ONSDigital/blaise-mi-extractcsv/storage/mysql"
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

	var found bool

	if encryptDestination, found = os.LookupEnv(util.EncryptLocation); !found {
		log.Fatal().Msg("The " + util.EncryptLocation + " variable has not been set")
		os.Exit(1)
	}

	gcloud := google.NewStorage()

	db := connect()

	service = extractor.NewService(&gcloud, db)
}

// sets up the database connection options and connects
func connect() *mysql.Storage {
	database, found := os.LookupEnv(util.Database)
	if !found {
		database = util.DefaultDatabase
	}

	server := func(db *mysql.Storage) {
		var s string
		if s, found = os.LookupEnv(util.Server); !found {
			log.Fatal().Msg("The " + util.Server + " varible has not been set")
			os.Exit(1)
		}
		db.Server = s
	}

	user := func(db *mysql.Storage) {
		var user string
		if user, found = os.LookupEnv(util.User); !found {
			log.Fatal().Msg("The " + util.Server + " varible has not been set")
			os.Exit(1)
		}
		db.User = user
	}

	password := func(db *mysql.Storage) {
		var pwd string
		if pwd, found = os.LookupEnv(util.Password); !found {
			log.Fatal().Msg("The " + util.Password + " varible has not been set")
			os.Exit(1)
		}
		db.Password = pwd
	}

	db := mysql.NewStorage(database, server, user, password)

	if err := db.Connect(); err != nil {
		// errors have already been reported and we can't continue
		os.Exit(1)
	}

	return db
}

// handle extract request events from publish / subscribe  queue
func ExtractFunction(_ context.Context, m PubSubMessage) error {
	// add additional actions as needed
	switch m.Action {
	case "extract_mi":
		return extractMi(m.Instrument)
	default:
		log.Warn().Msgf("message rejected, unknown action -> [%s]", m.Action)
		return nil
	}
}

func extractMi(instrument string) error {
	log.Info().Msgf("received extract_mi extract request for %s", instrument)

	var err error

	destination := instrument + ".csv"
	if err = service.ExtractMiInstrument(instrument, encryptDestination, destination); err != nil {
		return err
	}

	return nil
}
