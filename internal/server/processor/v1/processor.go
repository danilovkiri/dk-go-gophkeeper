package processor

import (
	"dk-go-gophkeeper/internal/server/cipher"
	"dk-go-gophkeeper/internal/server/processor"
	"dk-go-gophkeeper/internal/server/storage"
)

var (
	_ processor.Processor = (*Processor)(nil)
)

type Processor struct {
	storage   storage.DataStorage
	secretary cipher.Cipher
}

func InitService(st storage.DataStorage, cp cipher.Cipher) (*Processor, error) {
	serviceProcessor := &Processor{
		storage:   st,
		secretary: cp,
	}
	return serviceProcessor, nil
}
