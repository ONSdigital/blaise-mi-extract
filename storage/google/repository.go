package google

import (
	"archive/zip"
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
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

func (gs Storage) Zip(fileName, fromDirectory, toDirectory string) error {

	readBucket := gs.client.Bucket(fromDirectory)
	readObj := readBucket.Object(fileName)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	storageReader, err := readObj.NewReader(ctx)
	if err != nil {
		log.Err(err).Msgf("GCloud error: cannot create a reader")
		return err
	}

	defer func() { _ = storageReader.Close() }()

	writeBucket := gs.client.Bucket(toDirectory)
	writeObj := writeBucket.Object(fileName + ".zip")

	storageWriter := writeObj.NewWriter(ctx)
	defer func() { _ = storageWriter.Close() }()

	zipWriter := zip.NewWriter(storageWriter)
	defer func() { _ = zipWriter.Close() }()

	// add filename to zip
	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		log.Err(err).Msgf("error adding file to zip: %s in directory %s", fileName+".zip", toDirectory)
		return err
	}

	// copy from storage reader to zip writer
	_, err = io.Copy(zipFile, storageReader)

	if err != nil {
		log.Err(err).Msgf("error creating zip file: %s in directory %s", fileName+".zip", toDirectory)
		return err
	}

	log.Debug().Msgf("zip file: %s created in directory: %s", fileName+".zip", toDirectory)

	return nil
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

	log.Debug().Msgf("deleted file %s fromm directory: %s", file, directory)
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

	return encrypt([]*openpgp.Entity{recipient}, nil, storageReader, storageWriter)
}
