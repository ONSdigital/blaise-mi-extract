package google

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"os"
)

type Storage struct {
	client *storage.Client
}

func NewStorage() Storage {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Err(err).Msg("Cannot get GCloud Storage Bucket")
		os.Exit(1)
	}

	return Storage{client: client}
}
