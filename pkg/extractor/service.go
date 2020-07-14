package extractor

import (
	"encoding/csv"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
)

type Service interface {
	ExtractMiInstrument(instrument, destination, destinationFile string) (string, error)
}

type FileRepository interface {
	CreateFile(location, destinationFile string) (io.Writer, error)
}

type DBRepository interface {
	ExtractMIHeader(instrument string) (Instrument, error)
}

type Instrument struct {
	Spec string `db:"mi_spec"`
}

type MiSpec struct {
	SerialNumber string `json:"serial_number"`
	Hout         string `json:"hout"`
}

type service struct {
	fileRepository FileRepository
	dbRepository   DBRepository
}

func NewService(fileRepository FileRepository, dbRepository DBRepository) Service {
	return &service{fileRepository: fileRepository, dbRepository: dbRepository}
}

// extract data from the database and save as a csv
func (s service) ExtractMiInstrument(instrument, destination, destinationFile string) (string, error) {

	var headerJSON Instrument
	var err error

	if headerJSON, err = s.dbRepository.ExtractMIHeader(instrument); err != nil {
		return "", err // error already shown
	}

	// extract structure from the json
	var miSpec MiSpec
	err = json.Unmarshal([]byte(headerJSON.Spec), &miSpec)
	if err != nil {
		log.Warn().Msgf("cannot convert MiSpec to structure. Check structure definition")
		return "", err
	}

	// create the csv output file and give us a reference to the io.Writer so we can stream to it
	var c io.Writer

	c, err = s.fileRepository.CreateFile(destination, destinationFile)
	if err != nil {
		log.Err(err).Msgf("cannot create CSV file")
		return "", err
	}

	csvFile := csv.NewWriter(c)
	defer csvFile.Flush()

	// write the header
	err = csvFile.Write([]string{miSpec.SerialNumber, miSpec.Hout})
	if err != nil {
		log.Err(err).Msgf("cannot write CSV header")
		return "", err
	}

	// retrieve the data from the database*******************************

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
