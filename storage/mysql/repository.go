package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
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

	return &s
}

// connect to the database. Options (database, user etc.) have been set in NewStorage
func (s *Storage) Connect() error {

	var dbURI string
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")

	if !isSet {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", s.User, s.Password, s.Server,
			s.Database)
	} else {
		dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", s.User, s.Password, socketDir,
			s.Server, s.Database)
	}

	log.Info().Str("Connection string", dbURI).Msg("Connecting to DB")

	settings, err := mysql.ParseURL(dbURI)

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
