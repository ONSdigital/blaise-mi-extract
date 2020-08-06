package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/pkg"

)

// Proxy function to the real function
func ExtractFunction(ctx context.Context, m pkg.PubSubMessage) error {
	return pkg.HandleExtractionRequest(ctx, m)
}
