package config

import (
	"io/fs"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests

func TestNewDefaultConfiguration(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("DATABASE_DSN", "some_dsn")
	_ = os.Setenv("SERVER_ADDRESS", "some_server_address")
	_ = os.Setenv("USER_KEY", "some_user_key")
	_ = os.Setenv("BEARER_KEY", "some_key")
	_ = os.Setenv("BANK_CARD_DB", "someBankCard")
	_ = os.Setenv("LOGIN_PASSWORD_DB", "someLoginPassword")
	_ = os.Setenv("TEXT_BINARY_DB", "someTextBinary")
	_ = os.Setenv("HANDLERS_TO", "1000")
	cfg := NewDefaultConfiguration()
	var a = ""
	var c = ""
	var d = ""
	err := cfg.assignValues(&a, &c, &d)
	if err != nil {
		log.Fatal(err)
	}
	expCfg := Config{
		ServerAddress:   "some_server_address",
		DatabaseDSN:     "some_dsn",
		UserKey:         "some_user_key",
		AuthBearerName:  "some_key",
		BankCardDB:      "someBankCard",
		LoginPasswordDB: "someLoginPassword",
		TextBinaryDB:    "someTextBinary",
		HandlersTO:      1000,
	}
	assert.Equal(t, &expCfg, cfg)
}

func TestConfig_ParseFlags(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("USER_KEY", "some_user_key")
	cfg := NewDefaultConfiguration()
	os.Args = []string{"test", "-a", ":8080", "-c", "config_test.json"}
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}
	expCfg := Config{
		ServerAddress:   ":8080",
		DatabaseDSN:     "json_database_dsn",
		UserKey:         "some_user_key",
		AuthBearerName:  "token",
		BankCardDB:      "bankCard",
		LoginPasswordDB: "loginPassword",
		TextBinaryDB:    "textBinary",
		HandlersTO:      500,
	}
	assert.Equal(t, &expCfg, cfg)
}

func TestConfig_parseAppConfigPathError(t *testing.T) {
	os.Clearenv()
	cfg := NewDefaultConfiguration()
	var a = ""
	var c = "nonexistent_file.json"
	var d = ""
	err := cfg.assignValues(&a, &c, &d)
	var error *fs.PathError
	assert.ErrorAs(t, err, &error)
}
