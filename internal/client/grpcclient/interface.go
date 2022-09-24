package grpcclient

import "dk-go-gophkeeper/internal/client/storage/modelstroage"

// TextsBinariesGetter defines a set of methods for types implementing TextsBinariesGetter.
type TextsBinariesGetter interface {
	GetTextsBinaries() (map[string]modelstroage.TextOrBinary, error)
}

// LoginsPasswordsGetter defines a set of methods for types implementing LoginsPasswordsGetter.
type LoginsPasswordsGetter interface {
	GetLoginsPasswords() (map[string]modelstroage.LoginAndPassword, error)
}

// BankCardsGetter defines a set of methods for types implementing BankCardsGetter.
type BankCardsGetter interface {
	GetBankCards() (map[string]modelstroage.BankCard, error)
}

// BankCardSender defines a set of methods for types implementing BankCardSender.
type BankCardSender interface {
	SendBankCard(modelstroage.BankCard) error
}

// LoginPasswordSender defines a set of methods for types implementing LoginPasswordSender.
type LoginPasswordSender interface {
	SendLoginPassword(modelstroage.LoginAndPassword) error
}

// TextBinarySender defines a set of methods for types implementing TextBinarySender.
type TextBinarySender interface {
	SendTextBinary(modelstroage.TextOrBinary) error
}

// GRPCClient defines a set of embedded interfaces for types implementing GRPCClient.
type GRPCClient interface {
	TextsBinariesGetter
	LoginsPasswordsGetter
	BankCardsGetter
	BankCardSender
	LoginPasswordSender
	TextBinarySender
}
