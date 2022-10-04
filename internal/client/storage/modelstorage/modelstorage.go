// Package modelstorage provides models for local client data storage.
package modelstorage

type (
	LoginAndPassword struct {
		Identifier string
		Login      string
		Password   string
		Meta       string
	}
	TextOrBinary struct {
		Identifier string
		Entry      string
		Meta       string
	}
	BankCard struct {
		Identifier string
		Number     string
		Holder     string
		Cvv        string
		Meta       string
	}
	RegisterLogin struct {
		Login    string
		Password string
	}
)
