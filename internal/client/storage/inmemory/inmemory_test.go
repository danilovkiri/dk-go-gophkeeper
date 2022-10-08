package inmemory

import (
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"dk-go-gophkeeper/internal/config"
	"dk-go-gophkeeper/internal/mocks"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestInitStorage(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	_ = InitStorage(&logger, client, cfg)
}

func TestStorage_Remove(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)
	_ = st.AddBankCard("id1", "", "", "", "")
	_ = st.AddLoginPassword("id2", "", "", "")
	_ = st.AddTextBinary("id3", "", "")

	err := st.Remove("", "generic_db")
	assert.Equal(t, "identifier cannot be empty in db generic_db", err.Error())

	err = st.Remove("non_empty_id", "generic_db")
	assert.Equal(t, "invalid db generic_db", err.Error())

	client.EXPECT().RemoveBankCard(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.Remove("id1", cfg.BankCardDB)
	assert.Equal(t, "generic_error", err.Error())
	client.EXPECT().RemoveBankCard(gomock.Any()).Return(codes.OK, nil)
	err = st.Remove("id1", cfg.BankCardDB)
	assert.Equal(t, nil, err)

	client.EXPECT().RemoveLoginPassword(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.Remove("id2", cfg.LoginPasswordDB)
	assert.Equal(t, "generic_error", err.Error())
	client.EXPECT().RemoveLoginPassword(gomock.Any()).Return(codes.OK, nil)
	err = st.Remove("id2", cfg.LoginPasswordDB)
	assert.Equal(t, nil, err)

	client.EXPECT().RemoveTextBinary(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.Remove("id3", cfg.TextBinaryDB)
	assert.Equal(t, "generic_error", err.Error())
	client.EXPECT().RemoveTextBinary(gomock.Any()).Return(codes.OK, nil)
	err = st.Remove("id3", cfg.TextBinaryDB)
	assert.Equal(t, nil, err)

	_, ok := st.bankCardDB["id1"]
	assert.Equal(t, false, ok)
	_, ok = st.loginPasswordDB["id2"]
	assert.Equal(t, false, ok)
	_, ok = st.textBinaryDB["id3"]
	assert.Equal(t, false, ok)
}

func TestStorage_Login(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)
	_ = st.AddBankCard("id1", "", "", "", "")
	_ = st.AddLoginPassword("id2", "", "", "")
	_ = st.AddTextBinary("id3", "", "")

	err := st.Login("", "")
	assert.Equal(t, "Login/Password fields cannot be empty", err.Error())

	client.EXPECT().Login(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.Login("generic_login", "generic_password")
	assert.Equal(t, "generic_error", err.Error())

	client.EXPECT().Login(gomock.Any()).Return(codes.OK, nil)
	err = st.Login("generic_login", "generic_password")
	assert.Equal(t, nil, err)

	_, ok := st.bankCardDB["id1"]
	assert.Equal(t, false, ok)
	_, ok = st.loginPasswordDB["id2"]
	assert.Equal(t, false, ok)
	_, ok = st.textBinaryDB["id3"]
	assert.Equal(t, false, ok)
}

func TestStorage_Register(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)
	_ = st.AddBankCard("id1", "", "", "", "")
	_ = st.AddLoginPassword("id2", "", "", "")
	_ = st.AddTextBinary("id3", "", "")

	err := st.Register("", "")
	assert.Equal(t, "Login/Password fields cannot be empty", err.Error())

	client.EXPECT().Register(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.Register("generic_login", "generic_password")
	assert.Equal(t, "generic_error", err.Error())

	client.EXPECT().Register(gomock.Any()).Return(codes.OK, nil)
	err = st.Register("generic_login", "generic_password")
	assert.Equal(t, nil, err)

	_, ok := st.bankCardDB["id1"]
	assert.Equal(t, false, ok)
	_, ok = st.loginPasswordDB["id2"]
	assert.Equal(t, false, ok)
	_, ok = st.textBinaryDB["id3"]
	assert.Equal(t, false, ok)
}

func TestStorage_AddBankCard(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)

	err := st.AddBankCard("", "", "", "", "")
	assert.Equal(t, "identifier cannot be empty", err.Error())

	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	err = st.AddBankCard("id1", "", "", "", "")
	assert.Equal(t, nil, err)

	err = st.AddBankCard("id1", "", "", "", "")
	assert.Equal(t, "entry of type 'Bank Card' with ID id1 already exists", err.Error())

	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.AddBankCard("id2", "", "", "", "")
	assert.Equal(t, "generic_error", err.Error())
}

func TestStorage_AddLoginPassword(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)

	err := st.AddLoginPassword("", "", "", "")
	assert.Equal(t, "identifier cannot be empty", err.Error())

	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	err = st.AddLoginPassword("id1", "", "", "")
	assert.Equal(t, nil, err)

	err = st.AddLoginPassword("id1", "", "", "")
	assert.Equal(t, "entry of type 'Login And Password' with ID id1 already exists", err.Error())

	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.AddLoginPassword("id2", "", "", "")
	assert.Equal(t, "generic_error", err.Error())
}

func TestStorage_AddTextBinary(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)

	err := st.AddTextBinary("", "", "")
	assert.Equal(t, "identifier cannot be empty", err.Error())

	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)
	err = st.AddTextBinary("id1", "", "")
	assert.Equal(t, nil, err)

	err = st.AddTextBinary("id1", "", "")
	assert.Equal(t, "entry of type 'Text Or Binary' with ID id1 already exists", err.Error())

	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.Unknown, errors.New("generic_error"))
	err = st.AddTextBinary("id2", "", "")
	assert.Equal(t, "generic_error", err.Error())
}

func TestStorage_Get(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)
	_ = st.AddBankCard("id1", "", "", "", "")
	_ = st.AddLoginPassword("id2", "", "", "")
	_ = st.AddTextBinary("id3", "", "")

	_, err := st.Get("", "generic_db")
	assert.Equal(t, "identifier cannot be empty in db generic_db", err.Error())

	_, err = st.Get("non_empty_id", "generic_db")
	assert.Equal(t, "invalid db generic_db", err.Error())

	_, err = st.Get("id1", cfg.BankCardDB)
	assert.Equal(t, nil, err)
	_, err = st.Get("id2", cfg.LoginPasswordDB)
	assert.Equal(t, nil, err)
	_, err = st.Get("id3", cfg.TextBinaryDB)
	assert.Equal(t, nil, err)

	_, err = st.Get("nonexistent_id", cfg.BankCardDB)
	assert.Equal(t, "entry ID nonexistent_id in bankCard storage does not exist", err.Error())
	_, err = st.Get("nonexistent_id", cfg.LoginPasswordDB)
	assert.Equal(t, "entry ID nonexistent_id in loginPassword storage does not exist", err.Error())
	_, err = st.Get("nonexistent_id", cfg.TextBinaryDB)
	assert.Equal(t, "entry ID nonexistent_id in textBinary storage does not exist", err.Error())
}

func TestStorage_Sync(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.BankCardDB = "bankCard"
	cfg.LoginPasswordDB = "loginPassword"
	cfg.TextBinaryDB = "textBinary"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGRPCClient(ctrl)
	client.EXPECT().SendBankCard(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendLoginPassword(gomock.Any()).Return(codes.OK, nil)
	client.EXPECT().SendTextBinary(gomock.Any()).Return(codes.OK, nil)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	st := InitStorage(&logger, client, cfg)
	_ = st.AddBankCard("id1", "", "", "", "")
	_ = st.AddLoginPassword("id2", "", "", "")
	_ = st.AddTextBinary("id3", "", "")

	client.EXPECT().GetBankCards().Return(nil, codes.Unknown, errors.New("generic_error"))
	client.EXPECT().GetLoginsPasswords().Return(nil, codes.Unknown, errors.New("generic_error"))
	client.EXPECT().GetTextsBinaries().Return(nil, codes.Unknown, errors.New("generic_error"))
	err := st.Sync()
	assert.Equal(t, "generic_error", err.Error())

	cloudBankCardData := map[string]modelstorage.BankCard{"id4": {Identifier: "id4"}}
	cloudLoginPasswordData := map[string]modelstorage.LoginAndPassword{"id5": {Identifier: "id5"}}
	cloudTextBinaryData := map[string]modelstorage.TextOrBinary{"id6": {Identifier: "id6"}}
	client.EXPECT().GetBankCards().Return(cloudBankCardData, codes.OK, nil)
	client.EXPECT().GetLoginsPasswords().Return(cloudLoginPasswordData, codes.OK, nil)
	client.EXPECT().GetTextsBinaries().Return(cloudTextBinaryData, codes.OK, nil)
	err = st.Sync()
	assert.Equal(t, nil, err)

	_, ok := st.bankCardDB["id4"]
	assert.Equal(t, true, ok)
	_, ok = st.loginPasswordDB["id5"]
	assert.Equal(t, true, ok)
	_, ok = st.textBinaryDB["id6"]
	assert.Equal(t, true, ok)
}
