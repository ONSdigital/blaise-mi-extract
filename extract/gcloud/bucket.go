package gcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

var client *storage.Client

type GoogleBucket struct {
	Bucket string
}

func init() {
	var err error
	client, err = storage.NewClient(context.Background())
	if err != nil {
		log.Err(err).Str("internal Error", "Cannot get GCloud Storage Bucket")
		os.Exit(1)
	}
}

func (gb GoogleBucket) SaveToCSV(sourceFile, destinationFile string) error {

	log.Debug().Msgf("saving to GCloud Bucket; location: %s, sourceFile: %s, destinationFile: %s", gb.Bucket, sourceFile, destinationFile)

	ctx := context.Background()
	bh := client.Bucket(gb.Bucket)
	// Next check if the bucket exists
	if _, err := bh.Attrs(ctx); err != nil {
		return err
	}

	reader, err := os.Open(sourceFile)

	if err != nil {
		return err
	}

	defer func() { _ = reader.Close() }()

	obj := bh.Object(destinationFile)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	log.Debug().Msgf("file: %s, saved to: %s/%s", sourceFile, gb.Bucket, destinationFile)

	return nil
}
