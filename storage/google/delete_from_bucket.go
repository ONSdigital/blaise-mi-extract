package google

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

func (gs Storage) Delete(file, directory string) error {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := gs.client.Bucket(directory).Object(file)
	if err := o.Delete(ctx); err != nil {
		log.Warn().Msgf("delete of file %s fromm directory: %s failed", file, directory)
		return err
	}

	log.Debug().Msgf("file: %s/%s deleted", directory, file)

	return nil
}
