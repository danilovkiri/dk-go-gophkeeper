package modelstorage

type Removal struct {
	UserID     string
	Identifier string
	Db         string
}

type UserStorageEntry struct {
	ID           uint   `db:"id"`
	UserID       string `db:"user_id"`
	Login        string `db:"login"`
	Password     string `db:"password"`
	RegisteredAt string `db:"registered_at"`
}

type LoginPasswordStorageEntry struct {
	ID         uint   `db:"id"`
	UserID     string `db:"user_id"`
	Identifier string `db:"identifier"`
	Login      string `db:"login"`
	Password   string `db:"password"`
	Meta       string `db:"cred_meta"`
}

type BankCardStorageEntry struct {
	ID         uint   `db:"id"`
	UserID     string `db:"user_id"`
	Identifier string `db:"identifier"`
	Number     string `db:"card_number"`
	Holder     string `db:"card_holder"`
	CVV        string `db:"card_cvv"`
	Meta       string `db:"card_meta"`
}

type TextBinaryStorageEntry struct {
	ID         uint   `db:"id"`
	UserID     string `db:"user_id"`
	Identifier string `db:"identifier"`
	Entry      string `db:"text_entry"`
	Meta       string `db:"text_meta"`
}
