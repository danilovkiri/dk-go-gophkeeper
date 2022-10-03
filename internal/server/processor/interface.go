package processor

import (
	"context"
	"dk-go-gophkeeper/internal/server/modeldto"
)

type Authorizer interface {
	GetUserID(accessToken string) (string, error)
	GetTokenForUser(userID string) string
	AddNewUser(ctx context.Context, login, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

type Getter interface {
	GetBankCardData(ctx context.Context, userID, identifier string) ([]modeldto.BankCard, error)
	GetLoginPasswordData(ctx context.Context, userID, identifier string) ([]modeldto.LoginPassword, error)
	GetTextBinaryData(ctx context.Context, userID, identifier string) ([]modeldto.TextBinary, error)
}

type Setter interface {
	SetBankCardData(ctx context.Context, userID, identifier, number, holder, cvv, meta string) error
	SetLoginPasswordData(ctx context.Context, userID, identifier, login, password, meta string) error
	SetTextBinaryData(ctx context.Context, userID, identifier, entry, meta string) error
}

type Deleter interface {
	Delete(userID, identifier, db string)
}

type Processor interface {
	Authorizer
	Getter
	Setter
	Deleter
}
