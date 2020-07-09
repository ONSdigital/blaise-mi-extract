package extract

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/extract/gcloud"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var fpOnce sync.Once
var impl FilePersistence
var implName string

func googleImplementation() (FilePersistence, string) {

	var found bool
	var bucketName string

	if bucketName, found = os.LookupEnv(BucketKey); !found {
		log.Fatal().Msg("The " + BucketKey + " varible has not been set")
		os.Exit(1)
	}

	return &gcloud.GoogleBucket{Bucket: bucketName}, "Google"
}

func GetDefaultFilePersistenceImpl() FilePersistence {
	fpOnce.Do(func() { impl, implName = googleImplementation() })

	log.Debug().Msgf("returning %s FilePersistence implementation", implName)
	return impl
}

type FilePersistence interface {
	SaveToCSV(sourceFile, destinationFile string) error
}
