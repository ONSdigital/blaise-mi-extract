package blaise_mi_extractcsv

import "context"

type EncryptionMessage struct {
	FileName string `json:"fileName"`
	Bucket   string `json:"bucket"`
}

func ExtractMI(ctx context.Context, m EncryptionMessage) error {

	return nil
}
