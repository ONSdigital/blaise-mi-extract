package compressor

import (
	"io"
)

type Repository interface {
	CreateFile(location, destinationFile string) (io.Writer, error)
	DeleteFile(file, directory string) error
	CloseFile()
	ZipFile(fileName, fromDirectory, toDirectory string) (string, error)
}

type Service interface {
	ZipFile(fileName, fromDirectory, toDirectory string) (string, error)
	DeleteFile(file, directory string) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) DeleteFile(file, directory string) error {
	return s.r.DeleteFile(file, directory)
}

func (s service) ZipFile(fileName, fromDirectory, toDirectory string) (string, error) {
	return s.r.ZipFile(fileName, fromDirectory, toDirectory)
}
