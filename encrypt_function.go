package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/encrypt"
	"github.com/ONSDigital/blaise-mi-extract/util"
)

// handles event from item arriving in the encrypt  bucket
func EncryptFunction(_ context.Context, e util.GCSEvent) error {
	return encrypt.HandleEncryptionRequest(e.Name, e.Bucket)
}
