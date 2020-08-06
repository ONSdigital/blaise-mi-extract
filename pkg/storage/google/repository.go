package google

import (
	"archive/zip"
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Storage struct {
	client *storage.Client
	writer *storage.Writer
	ctx    context.Context
}

func NewStorage(ctx context.Context) Storage {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Err(err).Msg("Cannot get GCloud Storage Bucket")
		os.Exit(1)
	}

	return Storage{ctx: ctx, client: client}
}

func (gs *Storage) CreateFile(location, destinationFile string) (io.Writer, error) {

	log.Debug().Msgf("creating %s/%s", location, destinationFile)

	bh := gs.client.Bucket(location)
	// Next check if the bucket exists
	if _, err := bh.Attrs(gs.ctx); err != nil {
		return nil, err
	}

	obj := bh.Object(destinationFile)

	gs.writer = obj.NewWriter(gs.ctx)

	log.Debug().Msgf("file %s/%s created", location, destinationFile)

	return gs.writer, nil
}

func (gs *Storage) CloseFile() {
	if gs.writer != nil {
		err := gs.writer.Close()
		if err != nil {
			log.Err(err).Msg("close bucket writer failed")
			return
		}
		log.Debug().Msg("closed bucket writer")
	}
}

func (gs *Storage) DeleteFile(file, directory string) error {

	ctx, cancel := context.WithTimeout(gs.ctx, time.Second*10)
	defer cancel()

	o := gs.client.Bucket(directory).Object(file)
	if err := o.Delete(ctx); err != nil {
		log.Warn().Msgf("delete of file %s fromm directory: %s failed", file, directory)
		return err
	}

	log.Debug().Msgf("file: %s/%s deleted", directory, file)

	return nil
}

func (gs Storage) ZipFile(fileName, fromDirectory, toDirectory string) (string, error) {

	readBucket := gs.client.Bucket(fromDirectory)
	readObj := readBucket.Object(fileName)

	ctx, cancel := context.WithTimeout(gs.ctx, time.Second*60)
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

func (gs Storage) EncryptFile(publicKey, fileName, fromDirectory, toDirectory string) error {
	readBucket := gs.client.Bucket(fromDirectory)
	readObj := readBucket.Object(fileName)

	ctx, cancel := context.WithTimeout(gs.ctx, time.Second*15)
	defer cancel()

	storageReader, err := readObj.NewReader(ctx)
	if err != nil {
		log.Err(err).Msgf("cannot access %s/%s", fromDirectory, fileName)
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
		log.Err(err).Msgf("encrypt failed")
		return err
	}

	log.Info().Msgf("file %s encrypted and saved to %s/%s", fileName, toDirectory, fileName+".gpg")

	return nil
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}

	defer func() { _ = wc.Close() }()
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}

	return nil
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}

	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
