package storage

import (
	"context"
)

// BankCardAdder defines a set of methods for types implementing BankCardAdder.
type BankCardAdder interface {
	AddBankCard(identifier, number, holder, cvv, meta string) error
}

// LoginPasswordAdder defines a set of methods for types implementing LoginPasswordAdder.
type LoginPasswordAdder interface {
	AddLoginPassword(identifier, login, password, meta string) error
}

// TextBinaryAdder defines a set of methods for types implementing TextBinaryAdder.
type TextBinaryAdder interface {
	AddTextBinary(identifier, entry, meta string) error
}

// AllDataGetter defines a set of methods for types implementing AllDataGetter.
type AllDataGetter interface {
	ShowAllData() string
}

// Syncer defines a set of methods for types implementing Syncer.
type Syncer interface {
	Sync(ctx context.Context) error
}

// Remover defines a set of methods for types implementing Remover.
type Remover interface {
	Remove(string) (string, error)
}

// Cleaner defines a set of methods for types implementing Cleaner.
type Cleaner interface {
	CleanDB()
}

// Authorizer defines a set of methods for types implementing Authorizer.
type Authorizer interface {
	LoginRegister(login, password string) error
}

// DataStorage defines a set of embedded interfaces for types implementing DataStorage.
type DataStorage interface {
	BankCardAdder
	LoginPasswordAdder
	TextBinaryAdder
	Syncer
	AllDataGetter
	Remover
	Cleaner
	Authorizer
}
