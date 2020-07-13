package extractor

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

type Repository interface {
	Save(location, sourceFile, destinationFile string) error
}

type Service interface {
	DataToCSV() (string, error)
	SaveFile(location, sourceFile, destinationFile string) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) SaveFile(location, sourceFile, destinationFile string) error {
	return s.r.Save(location, sourceFile, destinationFile)
}

func (s service) DataToCSV() (string, error) {

	// get the data out of the database and save to a csv file

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
