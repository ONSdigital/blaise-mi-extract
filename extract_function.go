package blaise_mi_extractcsv

import (
	"context"
	"encoding/csv"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type PubSubMessage struct {
	Action     string `json:"action"`
	Instrument string `json:"instrument_name"`
}

const zipLocation = "ZIP_LOCATION"

type Extract struct {
	persistence Persistence
}

var extract Extract
var zipDestination string

func init() {
	util.Initialise()
	extract = Extract{persistence: GetPersistence()}
	var found bool

	if zipDestination, found = os.LookupEnv(zipLocation); !found {
		log.Fatal().Msg("The " + zipLocation + " varible has not been set for the google zipStorage provider")
		os.Exit(1)
	}

}

func ExtractFunction(ctx context.Context, m PubSubMessage) error {

	log.Info().
		Str("action", m.Action).
		Str("instrument", m.Instrument).
		Msgf("received request")

	if m.Action != "extract_mi" {
		log.Warn().Msgf("message rejected, unknown action -> [%s]", m.Action)
		return nil
	}

	var source string
	var err error

	if source, err = dataToCSV(); err != nil {
		return err
	}

	destination := m.Instrument + ".csv"

	if err = extract.persistence.Save(zipDestination, source, destination); err != nil {
		return err
	}

	return nil
}

func dataToCSV() (string, error) {

	// get the data out of the database and save to a csv file. optionally a zip file

	tmpFile, err := ioutil.TempFile("/tmp", "csv")
	if err != nil {
		log.Err(err).Msg("cannot create temporary file")
		return "", nil
	}

	defer func() { _ = tmpFile.Close() }()

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	var data = [][]string{{"Line1", "Hello Readers of"}, {"Line2", "golangcode.com"}}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Err(err).Msg("cannot write to temporary file")
			return "", err
		}
	}

	return tmpFile.Name(), nil
}
