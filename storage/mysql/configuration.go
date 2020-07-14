package mysql

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
)

func (s *Storage) loadConfiguration() {
	var found bool

	if s.server, found = os.LookupEnv(util.Server); !found {
		log.Fatal().Msg("The " + util.Server + " varible has not been set")
		os.Exit(1)
	}

	if s.database, found = os.LookupEnv(util.Database); !found {
		log.Fatal().Msg("The " + util.Database + " varible has not been set")
		os.Exit(1)
	}

	if s.user, found = os.LookupEnv(util.User); !found {
		log.Fatal().Msg("The " + util.Server + " varible has not been set")
		os.Exit(1)
	}

	if s.password, found = os.LookupEnv(util.Password); !found {
		log.Fatal().Msg("The " + util.Password + " varible has not been set")
		os.Exit(1)
	}

}
