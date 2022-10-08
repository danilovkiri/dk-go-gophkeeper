// Package cipher provides data ciphering functionality.
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"dk-go-gophkeeper/internal/config"
	procCipher "dk-go-gophkeeper/internal/server/cipher"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// check for interface compliance.
var (
	_ procCipher.Cipher = (*Cipher)(nil)
)

// Cipher defines attributes and methods of a Cipher instance.
type Cipher struct {
	aesgcm cipher.AEAD
	nonce  []byte
	key    []byte
	logger *zerolog.Logger
}

// NewCipherService initializes a Cipher instance.
func NewCipherService(cfg *config.Config, logger *zerolog.Logger) (*Cipher, error) {
	logger.Info().Msg("Attempting to initialize cipher")
	key := sha256.Sum256([]byte(cfg.UserKey))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	nonce := key[len(key)-aesgcm.NonceSize():]
	return &Cipher{
		aesgcm: aesgcm,
		nonce:  nonce,
		key:    []byte(cfg.UserKey),
		logger: logger,
	}, nil
}

// Encode performs ciphering data.
func (s *Cipher) Encode(data string) string {
	encoded := s.aesgcm.Seal(nil, s.nonce, []byte(data), nil)
	return hex.EncodeToString(encoded)
}

// Decode performs deciphering data.
func (s *Cipher) Decode(msg string) (string, error) {
	msgBytes, err := hex.DecodeString(msg)
	if err != nil {
		return "", err
	}
	decoded, err := s.aesgcm.Open(nil, s.nonce, msgBytes, nil)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// NewToken creates a new pair of user ID and its corresponding ciphered token.
func (s *Cipher) NewToken() (string, string) {
	userID := uuid.New().String()
	token := s.Encode(userID)
	return token, userID
}

// ValidateToken validates ciphered token and returns its corresponding user ID.
func (s *Cipher) ValidateToken(token string) (string, error) {
	userID, err := s.Decode(token)
	if err != nil {
		return "", err
	}
	return userID, nil
}
