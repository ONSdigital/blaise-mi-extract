package blaise_mi_extractcsv

import "github.com/ONSDigital/blaise-mi-extractcsv/persistence/google"

type Persistence interface {
	Zip(fileName, fromDirectory, toDirectory string) error
	Delete(file, directory string) error
	Save(location, sourceFile, destinationFile string) error
}

func GetPersistence() Persistence {
	return google.New()
}
