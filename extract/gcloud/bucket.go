package gcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"sync"
)

var client *storage.Client
var clientOnce sync.Once

type GoogleBucket struct{}

func (GoogleBucket) SaveToCSV(location, sourceFile, destinationFile string) error {

	clientOnce.Do(func() {
		// Pre-declare an err variable to avoid shadowing client.
		var err error
		client, err = storage.NewClient(context.Background())
		if err != nil {
			log.Err(err).
				Str("internal Error", "Cannot get GCloud Storage Bucket").
				Msgf("location: %s, sourceFile: %s, destinationFile: %s", location, sourceFile, destinationFile)
			return
		}
	})

	log.Info().
		Msgf("saving to GCloud Bucket; location: %s, sourceFile: %s, destinationFile: %s", location, sourceFile, destinationFile)

	return nil
}
