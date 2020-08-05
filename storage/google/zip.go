package google

import (
	"archive/zip"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"path/filepath"
	"strings"
	"time"
)

// compress a file and places in the compress location
func (gs Storage) Zip(fileName, fromDirectory, toDirectory string) (string, error) {

	readBucket := gs.client.Bucket(fromDirectory)
	readObj := readBucket.Object(fileName)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	storageReader, err := readObj.NewReader(ctx)
	if err != nil {
		log.Err(err).Msgf("GCloud error: cannot create a reader")
		return "", err
	}

	defer func() { _ = storageReader.Close() }()

	writeBucket := gs.client.Bucket(toDirectory)

	currentTime := time.Now()
	t := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // strip off .gpg suffix
	u := strings.TrimSuffix(t, filepath.Ext(t))               // strip off .csv suffix
	name := "mi_" + u + "_" + currentTime.Format("02012006") + "_" + currentTime.Format("150505") + ".compress"

	writeObj := writeBucket.Object(name)

	storageWriter := writeObj.NewWriter(ctx)
	defer func() { _ = storageWriter.Close() }()

	zipWriter := zip.NewWriter(storageWriter)
	defer func() { _ = zipWriter.Close() }()

	// add filename to compress
	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		log.Err(err).Msgf("error adding file to compress: %s in directory %s", name+".compress", toDirectory)
		return "", err
	}

	_, err = io.Copy(zipFile, storageReader)

	if err != nil {
		log.Err(err).Msgf("error creating compress file: %s in directory %s", fileName+".compress", toDirectory)
		return "", err
	}

	log.Debug().Msgf("file: %s, saved to: %s/%s", u, toDirectory, u)

	return name, nil
}
