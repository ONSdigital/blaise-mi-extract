package blaise_mi_extractcsv

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extractcsv/util"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

const encryptLocation = "ENCRYPT_LOCATION"

type Zip struct {
	persistence Persistence
}

var zip Zip
var encryptionDestination string

func init() {
	util.Initialise()
	zip = Zip{persistence: GetPersistence()}
	var found bool

	if encryptionDestination, found = os.LookupEnv(encryptLocation); !found {
		log.Fatal().Msg("The " + encryptLocation + " varible has not been set for the google zipStorage provider")
		os.Exit(1)
	}
}

func ZipFunction(ctx context.Context, e GCSEvent) error {

	log.Info().
		Str("bucket", e.Bucket).
		Str("file", e.Name).
		Msgf("received zip request")

	if err := zip.persistence.Zip(e.Name, e.Bucket, encryptionDestination); err != nil {
		log.Err(err).Msg("create zip failed")
		return err
	}

	if err := zip.persistence.Delete(e.Name, e.Bucket); err != nil {
		return err
	}

	return nil
}
