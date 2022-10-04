package storage

import (
	"context"
	"dk-go-gophkeeper/internal/server/storage/modelstorage"
)

type BatchDeleter interface {
	DeleteBatch(ctx context.Context, identifiers []string, userID, db string) error
	SendToQueue(item modelstorage.Removal)
	Flush(ctx context.Context, batch []modelstorage.Removal) error
}

type RegisterLogin interface {
	AddNewUser(ctx context.Context, login, password, userID string) error
	CheckUser(ctx context.Context, login, password string) (string, error)
}

type Getter interface {
	GetBankCardData(ctx context.Context, userID string) ([]modelstorage.BankCardStorageEntry, error)
	GetLoginPasswordData(ctx context.Context, userID string) ([]modelstorage.LoginPasswordStorageEntry, error)
	GetTextBinaryData(ctx context.Context, userID string) ([]modelstorage.TextBinaryStorageEntry, error)
}

type Setter interface {
	SetBankCardData(ctx context.Context, userID, identifier, number, holder, cvv, meta string) error
	SetLoginPasswordData(ctx context.Context, userID, identifier, login, password, meta string) error
	SetTextBinaryData(ctx context.Context, userID, identifier, entry, meta string) error
}

type DataStorage interface {
	RegisterLogin
	BatchDeleter
	Getter
	Setter
}
