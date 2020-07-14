package google

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"time"
)

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
