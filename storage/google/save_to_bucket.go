package google

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

func (gs Storage) Save(location, sourceFile, destinationFile string) error {

	log.Debug().Msgf("saving to GCloud Bucket; location: %s, sourceFile: %s, destinationFile: %s", location, sourceFile, destinationFile)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bh := gs.client.Bucket(location)
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

	log.Debug().Msgf("file: %s, saved to: %s/%s", sourceFile, location, destinationFile)

	return nil
}
