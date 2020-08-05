package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/extract"
)

// Proxy function to the real function
func ExtractFunction(_ context.Context, m extract.PubSubMessage) error {
	return extract.HandleExtractionRequest(m)
}
