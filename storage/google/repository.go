package google

import (
	"archive/zip"
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	name := "mi_" + u + "_" + currentTime.Format("02012006") + "_" + currentTime.Format("150505") + ".zip"

	writeObj := writeBucket.Object(name)

	storageWriter := writeObj.NewWriter(ctx)
	defer func() { _ = storageWriter.Close() }()

	zipWriter := zip.NewWriter(storageWriter)
	defer func() { _ = zipWriter.Close() }()

	// add filename to zip
	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		log.Err(err).Msgf("error adding file to zip: %s in directory %s", name+".zip", toDirectory)
		return "", err
	}

	_, err = io.Copy(zipFile, storageReader)

	if err != nil {
		log.Err(err).Msgf("error creating zip file: %s in directory %s", fileName+".zip", toDirectory)
		return "", err
	}

	log.Debug().Msgf("file: %s, saved to: %s/%s", u, toDirectory, u)

	return name, nil
}

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

func (gs Storage) Encrypt(publicKey, fileName, fromDirectory, toDirectory string) error {
	readBucket := gs.client.Bucket(fromDirectory)
	readObj := readBucket.Object(fileName)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	storageReader, err := readObj.NewReader(ctx)
	if err != nil {
		log.Err(err).Msgf("GCloud error: cannot create a reader")
		return err
	}

	defer func() { _ = storageReader.Close() }()

	writeBucket := gs.client.Bucket(toDirectory)
	writeObj := writeBucket.Object(fileName + ".gpg")

	storageWriter := writeObj.NewWriter(ctx)

	defer func() { _ = storageWriter.Close() }()

	// Read public key
	recipient, err := readEntity(publicKey)
	if err != nil {
		log.Err(err).Msgf("cannot read public key")
		return err
	}

	if err := encrypt([]*openpgp.Entity{recipient}, nil, storageReader, storageWriter); err != nil {
		log.Err(err).Msgf("encrypt failes")
		return err
	}

	log.Info().Msgf("file %s encrypted and saved to %s/%s", fileName, toDirectory, fileName+".gpg")

	return nil

}
