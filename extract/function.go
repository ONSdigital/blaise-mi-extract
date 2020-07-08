package extract

import (
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

var bucketName string

func init() {
	var found bool

	if bucketName, found = os.LookupEnv(BucketKey); !found {
		log.Fatal().Msg("The " + BucketKey + " varible has not been set")
		os.Exit(1)
	}

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
}

func ExtractMI(ctx context.Context, m PubSubMessage) error {

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

	if source, err = createCVS(); err != nil {
		return err
	}

	destination := m.Instrument + ".csv"

	if err = GetDefaultFilePersistenceImpl().SaveToCSV(bucketName, source, destination); err != nil {
		return err
	}

	return nil
}

func createCVS() (string, error) {

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
