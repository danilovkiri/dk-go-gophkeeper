// Package cipher provides data ciphering functionality.
package cipher

// Cipher defines a set of methods for types implementing Cipher.
type Cipher interface {
	Encode(data string) string
	Decode(msg string) (string, error)
	NewToken() (string, string)
	ValidateToken(token string) (string, error)
}
