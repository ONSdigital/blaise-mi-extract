package mysql

import (
	"github.com/rs/zerolog/log"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type Storage struct {
	DB       sqlbuilder.Database
	Server   string
	Database string
	User     string
	Password string
}

func NewStorage(database string, options ...func(*Storage)) *Storage {
	s := Storage{}

	s.Database = database

	for _, option := range options {
		if option != nil {
			option(&s)
		}
	}

	if err := s.Connect(); err != nil {
		log.Err(err).Msg("Cannot connect to database")
	}

	return &s
}

// connect to the database. Options (database, user etc.) have been set in NewStorage
func (s *Storage) Connect() error {
	var settings = mysql.ConnectionURL{
		Database: s.Database,
		Host:     s.Server,
		User:     s.User,
		Password: s.Password,
	}

	log.Debug().
		Str("databaseName", s.Database).
		Msg("Connecting to database")

	sess, err := mysql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", s.Database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", s.Database).
		Msg("Connected to database")

	s.DB = sess

	return nil
}

func (s Storage) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}
