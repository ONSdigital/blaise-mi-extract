package mysql

import (
	"database/sql"
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/extractor"
	"github.com/rs/zerolog/log"
)

type ResponseData struct {
	ResponseData string `db:"response_data"`
}

func (s Storage) ExtractMIHeader(instrument string) (extractor.Instrument, error) {
	res := s.DB.Collection("instrument").Find().Where("name = ? and phase = ?", instrument, "live")
	var i extractor.Instrument
	if err := res.One(&i); err != nil {
		log.Warn().Msgf("no instruments found in the database for %s or database error", instrument)
		return extractor.Instrument{}, err
	}
	return i, nil
}

func (s Storage) LoadResponseData(name string) (*sql.Rows, error) {

	rows, err := s.DB.Query(
		`SELECT response_data from case_response cr, instrument i, blaise.case c
			where c.instrument_id = i.id and
			cr.case_id = c.id and i.name = ?`, name)

	if err != nil {
		log.Err(err).Msgf("no instruments found in response_data for %s or database error", name)
		return nil, err
	}

	return rows, nil
}
