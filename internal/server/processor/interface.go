// Package processor provides intermediary layer functionality between the DB and handlers.
package processor

import (
	"context"
	"dk-go-gophkeeper/internal/server/modeldto"
)

// Authorizer defines a set of methods for types implementing Authorizer.
type Authorizer interface {
	GetUserID(accessToken string) (string, error)
	AddNewUser(ctx context.Context, login, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

// Getter defines a set of methods for types implementing Getter.
type Getter interface {
	GetBankCardData(ctx context.Context, userID string) ([]modeldto.BankCard, error)
	GetLoginPasswordData(ctx context.Context, userID string) ([]modeldto.LoginPassword, error)
	GetTextBinaryData(ctx context.Context, userID string) ([]modeldto.TextBinary, error)
}

// Setter defines a set of methods for types implementing Setter.
type Setter interface {
	SetBankCardData(ctx context.Context, userID, identifier, number, holder, cvv, meta string) error
	SetLoginPasswordData(ctx context.Context, userID, identifier, login, password, meta string) error
	SetTextBinaryData(ctx context.Context, userID, identifier, entry, meta string) error
}

// Deleter defines a set of methods for types implementing Deleter.
type Deleter interface {
	Delete(userID, identifier, db string)
}

// Processor defines a set of methods for types implementing Processor.
type Processor interface {
	Authorizer
	Getter
	Setter
	Deleter
}
