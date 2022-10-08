package processor

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	"dk-go-gophkeeper/internal/mocks"
	"dk-go-gophkeeper/internal/server/modeldto"
	"dk-go-gophkeeper/internal/server/storage/modelstorage"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestInitService(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	_ = InitService(storage, cipher, log.New(os.Stdout, "test", 0))
}

func TestProcessor_GetUserID(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	cipher.EXPECT().ValidateToken(gomock.Any()).Return("generic_user_id", nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	userID, err := processor.GetUserID("generic_access_token")
	assert.Equal(t, "generic_user_id", userID)
	assert.Equal(t, nil, err)
}

func TestProcessor_AddNewUser(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	cipher.EXPECT().NewToken().Return("generic_access_token", "generic_user_id")
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	storage.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	accessToken, err := processor.AddNewUser(context.Background(), "generic_login", "generic_password")
	assert.Equal(t, "generic_access_token", accessToken)
	assert.Equal(t, nil, err)
}

func TestProcessor_LoginUser(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("generic_user_id", nil)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	accessToken, err := processor.LoginUser(context.Background(), "generic_login", "generic_password")
	assert.Equal(t, "generic_encoded_data", accessToken)
	assert.Equal(t, nil, err)
}

func TestProcessor_LoginUserFail(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("generic_error"))
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.LoginUser(context.Background(), "generic_login", "generic_password")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetBankCardData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("generic_decoded_data", nil).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.BankCardStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Number:     cipher.Encode("some_data"),
			Holder:     cipher.Encode("some_data"),
			CVV:        cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	bankCards, err := processor.GetBankCardData(context.Background(), "some_user_id")
	assert.Equal(t, nil, err)
	expectedBankCards := []modeldto.BankCard{{Identifier: "generic_decoded_data", Number: "generic_decoded_data", Holder: "generic_decoded_data", CVV: "generic_decoded_data", Meta: "generic_decoded_data"}}
	assert.Equal(t, expectedBankCards, bankCards)
}

func TestProcessor_GetBankCardDataFail1(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetBankCardData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetBankCardDataFail2(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("", errors.New("generic_error")).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.BankCardStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Number:     cipher.Encode("some_data"),
			Holder:     cipher.Encode("some_data"),
			CVV:        cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetBankCardData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetLoginPasswordData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("generic_decoded_data", nil).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.LoginPasswordStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Login:      cipher.Encode("some_data"),
			Password:   cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	loginsPasswords, err := processor.GetLoginPasswordData(context.Background(), "some_user_id")
	assert.Equal(t, nil, err)
	expectedLoginsPasswords := []modeldto.LoginPassword{{Identifier: "generic_decoded_data", Login: "generic_decoded_data", Password: "generic_decoded_data", Meta: "generic_decoded_data"}}
	assert.Equal(t, expectedLoginsPasswords, loginsPasswords)
}

func TestProcessor_GetLoginPasswordDataFail1(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetLoginPasswordData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetLoginPasswordDataFail2(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("", errors.New("generic_error")).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.LoginPasswordStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Login:      cipher.Encode("some_data"),
			Password:   cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetLoginPasswordData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetTextBinaryData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("generic_decoded_data", nil).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.TextBinaryStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Entry:      cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	textsBinaries, err := processor.GetTextBinaryData(context.Background(), "some_user_id")
	assert.Equal(t, nil, err)
	expectedTextsBinaries := []modeldto.TextBinary{{Identifier: "generic_decoded_data", Entry: "generic_decoded_data", Meta: "generic_decoded_data"}}
	assert.Equal(t, expectedTextsBinaries, textsBinaries)
}

func TestProcessor_GetTextBinaryDataFail1(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetTextBinaryData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_GetTextBinaryDataFail2(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	cipher.EXPECT().Decode(gomock.Any()).Return("", errors.New("generic_error")).AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storageOutput := []modelstorage.TextBinaryStorageEntry{
		{
			ID:         1,
			UserID:     cipher.Encode("some_data"),
			Identifier: cipher.Encode("some_data"),
			Entry:      cipher.Encode("some_data"),
			Meta:       cipher.Encode("some_data"),
		},
	}
	storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(storageOutput, nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	_, err := processor.GetTextBinaryData(context.Background(), "some_user_id")
	assert.Equal(t, "generic_error", err.Error())
}

func TestProcessor_SetBankCardData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().SetBankCardData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	err := processor.SetBankCardData(context.Background(), "", "", "", "", "", "")
	assert.Equal(t, nil, err)
}

func TestProcessor_SetLoginPasswordData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().SetLoginPasswordData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	err := processor.SetLoginPasswordData(context.Background(), "", "", "", "", "")
	assert.Equal(t, nil, err)
}

func TestProcessor_SetTextBinaryData(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().SetTextBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	err := processor.SetTextBinaryData(context.Background(), "", "", "", "")
	assert.Equal(t, nil, err)
}

func TestProcessor_Delete(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cipher := mocks.NewMockCipher(ctrl)
	cipher.EXPECT().Encode(gomock.Any()).Return("generic_encoded_data").AnyTimes()
	storage := mocks.NewMockDataStorage(ctrl)
	storage.EXPECT().SendToQueue(gomock.Any()).Return()
	processor := InitService(storage, cipher, log.New(os.Stdout, "test", 0))
	processor.Delete("", "", "")
}
