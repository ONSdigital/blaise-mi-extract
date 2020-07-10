package persistence

import (
	"github.com/ONSDigital/blaise-mi-extractcsv/persistence/gcloud"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"sync"
)

var fpOnce sync.Once
var impl FilePersistence
var implName string

const ProviderKey string = "FILE_PROVIDER"

var defaultProvider = googleStorage

func googleStorage() (FilePersistence, string) {
	return gcloud.New()
}

func GetStorageProvider() FilePersistence {
	fpOnce.Do(func() { impl, implName = fileProvider()() })
	log.Debug().Msgf("returning %s FilePersistence implementation", implName)
	return impl
}

func fileProvider() func() (FilePersistence, string) {
	var found bool
	var provider string

	// add additional providers here
	if provider, found = os.LookupEnv(ProviderKey); found {
		switch strings.ToLower(provider) {
		case "google":
			return googleStorage
		}
	}
	// default provider
	return defaultProvider
}

type FilePersistence interface {
	Save(location, sourceFile, destinationFile string) error
	Zip(fileName, fromDirectory, toDirectory string) error
	Delete(file, directory string) error
}
