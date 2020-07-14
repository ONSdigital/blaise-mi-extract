package google

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
)

// create a file in a bucket and return the io.Writer for streaming

func (gs Storage) CreateFile(location, destinationFile string) (io.Writer, error) {

	log.Debug().Msgf("create %s/%s", location, destinationFile)

	ctx := context.Background()

	bh := gs.client.Bucket(location)
	// Next check if the bucket exists
	if _, err := bh.Attrs(ctx); err != nil {
		return nil, err
	}

	obj := bh.Object(destinationFile)
	return obj.NewWriter(ctx), nil
}
