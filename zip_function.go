package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/compress"
	"github.com/ONSDigital/blaise-mi-extract/util"
)

// Proxy function to the real function
func ZipFunction(_ context.Context, e util.GCSEvent) error {
	return compress.ZipCompress(e.Name, e.Bucket)
}
