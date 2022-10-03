package grpcclient

import (
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"google.golang.org/grpc/codes"
)

// TextsBinariesGetter defines a set of methods for types implementing TextsBinariesGetter.
type TextsBinariesGetter interface {
	GetTextsBinaries() (map[string]modelstorage.TextOrBinary, codes.Code, error)
}

// LoginsPasswordsGetter defines a set of methods for types implementing LoginsPasswordsGetter.
type LoginsPasswordsGetter interface {
	GetLoginsPasswords() (map[string]modelstorage.LoginAndPassword, codes.Code, error)
}

// BankCardsGetter defines a set of methods for types implementing BankCardsGetter.
type BankCardsGetter interface {
	GetBankCards() (map[string]modelstorage.BankCard, codes.Code, error)
}

// BankCardSender defines a set of methods for types implementing BankCardSender.
type BankCardSender interface {
	SendBankCard(modelstorage.BankCard) (codes.Code, error)
}

// LoginPasswordSender defines a set of methods for types implementing LoginPasswordSender.
type LoginPasswordSender interface {
	SendLoginPassword(modelstorage.LoginAndPassword) (codes.Code, error)
}

// TextBinarySender defines a set of methods for types implementing TextBinarySender.
type TextBinarySender interface {
	SendTextBinary(modelstorage.TextOrBinary) (codes.Code, error)
}

// Remover defines a set of methods for types implementing Remover.
type Remover interface {
	RemoveBankCard(string) (codes.Code, error)
	RemoveLoginPassword(string) (codes.Code, error)
	RemoveTextBinary(string) (codes.Code, error)
}

// Authorizer defines a set of methods for types implementing Authorizer.
type Authorizer interface {
	Login(modelstorage.RegisterLogin) (codes.Code, error)
	Register(modelstorage.RegisterLogin) (codes.Code, error)
}

// GRPCClient defines a set of embedded interfaces for types implementing GRPCClient.
type GRPCClient interface {
	TextsBinariesGetter
	LoginsPasswordsGetter
	BankCardsGetter
	BankCardSender
	LoginPasswordSender
	TextBinarySender
	Remover
	Authorizer
}
