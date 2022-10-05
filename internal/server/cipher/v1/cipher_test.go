package cipher

import (
	"dk-go-gophkeeper/internal/config"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

func TestCipher_NewToken(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cipher, _ := NewCipherService(cfg, log.New(os.Stdout, "test", 0))
	token, userID := cipher.NewToken()
	assert.Equal(t, cipher.Encode(userID), token)
}

func TestCipher_ValidateToken(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cipher, _ := NewCipherService(cfg, log.New(os.Stdout, "test", 0))
	token, expUserID := cipher.NewToken()
	obsUserID, _ := cipher.ValidateToken(token)
	assert.Equal(t, expUserID, obsUserID)
}

func TestCipher_ValidateTokenFail(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cipher, _ := NewCipherService(cfg, log.New(os.Stdout, "test", 0))
	_, err := cipher.ValidateToken("non-hex-encoded-data")
	assert.Equal(t, err.Error(), "encoding/hex: invalid byte: U+006E 'n'")
}

func TestDecode_Fail(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cipher, _ := NewCipherService(cfg, log.New(os.Stdout, "test", 0))
	var newNonce []byte
	for i := 0; i < len(cipher.nonce); i++ {
		newNonce = append(newNonce, 1)
	}
	cipher.nonce = newNonce
	res, err := cipher.Decode("c277fd4361e8c0e81e90bc030a31621ff6ef71503544154b7f0e29aae1f69dec0a00")
	if err != nil {
		assert.Equal(t, err.Error(), "cipher: message authentication failed")
	}
	assert.Equal(t, "", res)
}

type CipherTestSuite struct {
	suite.Suite
	cipher *Cipher
	config *config.Config
}

func (suite *CipherTestSuite) SetupTest() {
	suite.config = config.NewDefaultConfiguration()
	suite.config.UserKey = "jds__63h3_7ds"
	suite.cipher, _ = NewCipherService(suite.config, log.New(os.Stdout, "test", 0))
}

func TestSecretaryTestSuite(t *testing.T) {
	suite.Run(t, new(CipherTestSuite))
}

func (suite *CipherTestSuite) TestEncode() {
	tests := []struct {
		name             string
		data             string
		expectedEncoding string
	}{
		{
			name:             "sample 1",
			data:             "sample text string",
			expectedEncoding: "c277fd4361e8c0e81e90bc030a31621ff6ef71503544154b7f0e29aae1f69dec0a00",
		},
		{
			name:             "sample 2",
			data:             "another integer data piece",
			expectedEncoding: "d078ff4765e892bc1286bc461e206256fce9061c0fffc7ae409a76a2c8fd0933da10a997181b1f89e06e",
		},
	}

	// perform each test
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedEncoding, suite.cipher.Encode(tt.data))
		})
	}
}

func (suite *CipherTestSuite) TestDecode() {
	var invalidByteError *hex.InvalidByteError
	tests := []struct {
		name             string
		expectedDecoding string
		data             string
		error            error
	}{
		{
			name:             "sample 1",
			expectedDecoding: "sample text string",
			data:             "c277fd4361e8c0e81e90bc030a31621ff6ef71503544154b7f0e29aae1f69dec0a00",
			error:            nil,
		},
		{
			name:             "sample 2",
			expectedDecoding: "another integer data piece",
			data:             "d078ff4765e892bc1286bc461e206256fce9061c0fffc7ae409a76a2c8fd0933da10a997181b1f89e06e",
			error:            nil,
		},
		{
			name:             "sample 3",
			expectedDecoding: "",
			data:             "non-hex-encoded-data",
			error:            invalidByteError,
		},
		{
			name:             "sample 4",
			expectedDecoding: "",
			data:             "d078ff4765e892bc1286bc461e206256fce9061c0fffc7ae409a76a",
			error:            nil,
		},
	}

	// perform each test
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			res, err := suite.cipher.Decode(tt.data)
			if err != nil {
				assert.ErrorAs(t, err, &tt.error)
			}
			assert.Equal(t, tt.expectedDecoding, res)

		})
	}
}
