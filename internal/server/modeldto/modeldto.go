// Package modeldto provides models for data transferring between the handlers and the storage.
package modeldto

type LoginPassword struct {
	Identifier string
	Login      string
	Password   string
	Meta       string
}

type BankCard struct {
	Identifier string
	Number     string
	Holder     string
	CVV        string
	Meta       string
}

type TextBinary struct {
	Identifier string
	Entry      string
	Meta       string
}
