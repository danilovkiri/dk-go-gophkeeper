package grpcclient

import "dk-go-gophkeeper/internal/client/storage/modelstorage"

// TextsBinariesGetter defines a set of methods for types implementing TextsBinariesGetter.
type TextsBinariesGetter interface {
	GetTextsBinaries() (map[string]modelstorage.TextOrBinary, error)
}

// LoginsPasswordsGetter defines a set of methods for types implementing LoginsPasswordsGetter.
type LoginsPasswordsGetter interface {
	GetLoginsPasswords() (map[string]modelstorage.LoginAndPassword, error)
}

// BankCardsGetter defines a set of methods for types implementing BankCardsGetter.
type BankCardsGetter interface {
	GetBankCards() (map[string]modelstorage.BankCard, error)
}

// BankCardSender defines a set of methods for types implementing BankCardSender.
type BankCardSender interface {
	SendBankCard(modelstorage.BankCard) error
}

// LoginPasswordSender defines a set of methods for types implementing LoginPasswordSender.
type LoginPasswordSender interface {
	SendLoginPassword(modelstorage.LoginAndPassword) error
}

// TextBinarySender defines a set of methods for types implementing TextBinarySender.
type TextBinarySender interface {
	SendTextBinary(modelstorage.TextOrBinary) error
}

// Remover defines a set of methods for types implementing Remover.
type Remover interface {
	Remove(string) error
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
}
