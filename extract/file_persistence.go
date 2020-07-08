package extract

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/extract/gcloud"
	"github.com/rs/zerolog/log"
	"sync"
)

var fpOnce sync.Once

func googleImplementation() (FilePersistence, string) {
	return &gcloud.GoogleBucket{}, "Google"
}

func GetDefaultFilePersistenceImpl() FilePersistence {
	var impl FilePersistence
	var implName string

	// lazy initialisation
	fpOnce.Do(func() { impl, implName = googleImplementation() })

	log.Debug().Msgf("returning %s FilePersistence implementation", implName)
	return impl
}

type FilePersistence interface {
	SaveToCSV(location, sourceFile, destinationFile string) error
}
