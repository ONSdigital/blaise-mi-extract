package mysql

import (
	"github.com/rs/zerolog/log"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type Storage struct {
	DB       sqlbuilder.Database
	server   string
	database string
	user     string
	password string
}

func NewStorage() Storage {
	s := Storage{}
	s.loadConfiguration()

	if err := s.Connect(); err != nil {
		log.Err(err).Msg("Cannot connect to mysql")
	}

	return s
}

func (s *Storage) Connect() error {
	var settings = mysql.ConnectionURL{
		Database: s.database,
		Host:     s.server,
		User:     s.user,
		Password: s.password,
	}

	log.Debug().
		Str("databaseName", s.database).
		Msg("Connecting to database")

	sess, err := mysql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", s.database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", s.database).
		Msg("Connected to database")

	s.DB = sess

	return nil
}

func (s Storage) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}
