package storage

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

// Syncer defines a set of methods for types implementing Syncer.
type Syncer interface {
	Sync() error
}

// Remover defines a set of methods for types implementing Remover.
type Remover interface {
	Remove(string, string) error
}

// Getter defines a set of methods for types implementing Getter.
type Getter interface {
	Get(string, string) (string, error)
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
	Getter
	Syncer
	Remover
	Cleaner
	Authorizer
}
