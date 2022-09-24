package modelstroage

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
)
