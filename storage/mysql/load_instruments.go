package mysql

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/pkg/extractor"
	"github.com/rs/zerolog/log"
)

func (s Storage) GetMISpecs(instrument string) ([]extractor.MISpec, error) {

	var specs []extractor.MISpec

	q := s.DB.Select("header_name", "response_key").
		From("instrument", "mi_spec", "mi_values").
		Where("instrument.mi_spec = mi_spec.id").
		And("mi_spec.id = mi_values.spec_id").
		And("instrument.name = ?", instrument).
		And("instrument.phase = ?", "live")

	if err := q.All(&specs); err != nil {
		log.Warn().Msgf("no instruments found or no mi specs for %s or database error", instrument)
		return specs, err
	}
	return specs, nil
}

func (s Storage) LoadResponseData(name string) ([]extractor.ResponseData, error) {

	var responses []extractor.ResponseData

	q := s.DB.Select("response_data").
		From("case_response cr", "instrument i", "blaise.case c").
		Where("c.instrument_id = i.id").
		And("cr.case_id = c.id").
		And("i.name = ?", name)

	if err := q.All(&responses); err != nil {
		log.Warn().Msgf("no responses found for %s or database error", name)
		return responses, err
	}

	return responses, nil
}
