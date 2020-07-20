package extractor

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
)

type Service interface {
	ExtractMiInstrument(instrument, destination, destinationFile string) error
}

type FileRepository interface {
	CreateFile(location, destinationFile string) (io.Writer, error)
	CloseFile()
}

type DBRepository interface {
	ExtractMIHeader(instrument string) (Instrument, error)
	LoadResponseData(name string) (*sql.Rows, error)
}

type Instrument struct {
	Spec string `db:"mi_spec"`
}

type service struct {
	fileRepository FileRepository
	dbRepository   DBRepository
}

// create a new service instance
func NewService(fileRepository FileRepository, dbRepository DBRepository) Service {
	return &service{fileRepository: fileRepository, dbRepository: dbRepository}
}

// extract data from the database and save as a csv
func (s service) ExtractMiInstrument(instrument, destination, destinationFile string) error {

	var headerJSON Instrument
	var err error

	if headerJSON, err = s.dbRepository.ExtractMIHeader(instrument); err != nil {
		return err // error already shown
	}

	// extract structure from the json
	var miSpec = map[string]string{}
	err = json.Unmarshal([]byte(headerJSON.Spec), &miSpec)

	if err != nil {
		log.Warn().Msgf("cannot convert MiSpec to structure. Check structure definition")
		return err
	}

	// defer calling this until we know we actually have some data
	var createCSV = func() (*csv.Writer, error) {
		var c io.Writer

		c, err = s.fileRepository.CreateFile(destination, destinationFile)
		if err != nil {
			log.Err(err).Msgf("cannot create CSV file")
			return nil, err
		}
		csvFile := csv.NewWriter(c)
		return csvFile, nil
	}

	// defer calling this until we know we actually have data
	var writeHeader = func(csvFile *csv.Writer) error {
		// write the header
		keys := make([]string, 0, len(miSpec))
		for key := range miSpec {
			keys = append(keys, key)
		}
		err = csvFile.Write(keys)
		if err != nil {
			log.Err(err).Msgf("cannot write CSV header")
			return err
		}
		return nil
	}

	rows, err := s.dbRepository.LoadResponseData(instrument)
	if err != nil {
		log.Err(err).Msg("cannot load response data")
		return nil
	}

	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	var csvCreated = false
	var csvFile *csv.Writer

	for rows.Next() {
		// we have at least one row
		// so create the file and write the header; Ugly
		if !csvCreated {
			if csvFile, err = createCSV(); err != nil {
				log.Err(err).Msg("cannot create CSV file")
				return nil
			}
			csvCreated = true
			if err := writeHeader(csvFile); err != nil {
				log.Err(err).Msg("cannot write CSV header")
				return nil
			}
		}

		var js string
		err := rows.Scan(&js)
		if err != nil {
			log.Err(err).Msg("row scan failed")
			return nil
		}

		m := map[string]string{}
		err = json.Unmarshal([]byte(js), &m)
		if err != nil {
			log.Err(err).Msg("invalid json string in response_data")
			return nil
		}

		var r []string

		for _, v := range miSpec { // iterate over header values
			if val, ok := m[v]; ok {
				r = append(r, val)
			} else {
				r = append(r, "")
			}
		}

		err = csvFile.Write(r)
		if err != nil {
			log.Err(err).Msgf("cannot write CSV row")
			return err
		}

		r = nil

	}

	if !csvCreated {
		log.Warn().Msgf("no response data found for instrument: %s", instrument)
		return nil
	}

	csvFile.Flush()

	s.fileRepository.CloseFile()

	return nil
}
