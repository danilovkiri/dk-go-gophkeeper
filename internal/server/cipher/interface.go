package cipher

type Cipher interface {
	Encode(data string) string
	Decode(msg string) (string, error)
	NewToken() (string, string)
	ValidateToken(token string) (string, error)
}
