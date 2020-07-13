package encryption

type Repository interface {
	Encrypt(publicKey, fileName, fromDirectory, toDirectory string) error
	Delete(file, directory string) error
}

type Service interface {
	EncryptFile(publicKey, fileName, fromDirectory, toDirectory string) error
	DeleteFile(file, directory string) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) EncryptFile(publicKey, fileName, fromDirectory, toDirectory string) error {
	return s.r.Encrypt(publicKey, fileName, fromDirectory, toDirectory)
}

func (s service) DeleteFile(file, directory string) error {
	return s.r.Delete(file, directory)
}
