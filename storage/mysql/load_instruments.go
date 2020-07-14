package mysql

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/extractor"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"upper.io/db.v3/lib/sqlbuilder"
)

type ResponseData struct {
	ResponseData string `db:"response_data"`
}

// ATM we only have one instrument id
// We need one of these for every instrument type

func (s Storage) ExtractMIHeader(instrument string) (extractor.Instrument, error) {
	res := s.DB.Collection("instrument").Find().Where("name = ? and phase = ?", instrument, "live")
	var i extractor.Instrument
	if err := res.One(&i); err != nil {
		log.Warn().Msgf("no instruments found in the database for %s or database error", instrument)
		return extractor.Instrument{}, err
	}
	return i, nil
}

// return a function that, when called, retrieves the data row-by-row
func (s Storage) LoadResponseData(name string) error {

	rows, err := s.DB.Query(
		"SELECT response_data from case_response cr, instrument i , where cr.instrument_id = i.id and i.name = ?", name)

	if err != nil {
		log.Err(err).Msgf("no instruments found in response_data for %s or database error", name)
		return err
	}

	defer func() { _ = rows.Close() }()

	s.iter = sqlbuilder.NewIterator(rows)

	return nil
}

func (s Storage) Next() ResponseData {
	var responseData ResponseData
	eof := s.iter.Next(&responseData)
	if eof {
		_ = s.iter.Close()
	}
	return responseData
}
