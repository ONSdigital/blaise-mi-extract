package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/ONSDigital/blaise-mi-extractcsv"
	"log"
	"os"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func main() {
	funcframework.RegisterEventFunction("/", blaise_mi_extractcsv.ExtractFunction)
	funcframework.RegisterEventFunction("/zip", blaise_mi_extractcsv.ZipFunction)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
