package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	blaise-mi-extract "github.com/ONSDigital/blaise-mi-extract"
	"log"
	"os"
)

// emulates the cloud functions
func main() {
	funcframework.RegisterEventFunction("/extract", blaise-mi-extract.ExtractFunction)
	funcframework.RegisterEventFunction("/zip", blaise-mi-extract.ZipFunction)
	funcframework.RegisterEventFunction("/encrypt", blaise-mi-extract.EncryptFunction)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}

}
