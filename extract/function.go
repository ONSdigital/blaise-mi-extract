package extract

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"time"
)

type PubSubMessage struct {
	Action     string `json:"action"`
	Instrument string `json:"instrument_name"`
}

const (
	BucketKey string = "MI_BUCKET_NAME"
	LogFormat string = "LOG_FORMAT"
	Terminal         = "Terminal"
	Json             = "Json"
	Debug            = "DEBUG"
)

var persistFile FilePersistence

func init() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// change log format
	if terminal, isFound := os.LookupEnv(LogFormat); isFound {
		switch terminal {
		case Terminal:
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.Stamp})
		case Json:
			// json is the default
		}
	}

	if debug, f := os.LookupEnv(Debug); f {
		switch debug {
		case "True":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	}

	persistFile = GetDefaultFilePersistenceImpl()
}

func MIToCSV(ctx context.Context, m PubSubMessage) error {

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

	if source, err = dataToCSV(true); err != nil {
		return err
	}

	destination := m.Instrument + ".csv"

	if err = persistFile.SaveToCSV(source, destination); err != nil {
		return err
	}

	return nil
}

func dataToCSV(zip bool) (string, error) {

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
