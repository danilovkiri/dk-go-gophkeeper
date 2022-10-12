package grpcclient

import (
	"context"
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/mocks"
	"dk-go-gophkeeper/internal/server/api/handlers"
	"dk-go-gophkeeper/internal/server/api/interceptors"
	"dk-go-gophkeeper/internal/server/cipher/v1"
	serverStorage "dk-go-gophkeeper/internal/server/storage/modelstorage"
	"errors"
	"log"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type ClientTestSuite struct {
	suite.Suite
	storage *mocks.MockDataStorage
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	server  *handlers.GophkeeperServer
	s       *grpc.Server
	client  *GRPCClient
	cfg     *config.Config
	cipher  *cipher.Cipher
}

func (suite *ClientTestSuite) SetupTest() {
	cfg := config.NewDefaultConfiguration()
	cfg.ServerAddress = ":8080"
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	suite.cfg = cfg
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.wg = &sync.WaitGroup{}
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	suite.storage = mocks.NewMockDataStorage(ctrl)
	server, err := handlers.InitServer(cfg, suite.storage, &logger)
	if err != nil {
		log.Fatal(err)
	}
	suite.server = server
	cipherInstance, err := cipher.NewCipherService(cfg, &logger)
	suite.cipher = cipherInstance
	if err != nil {
		log.Fatal(err)
	}
	interceptorService := interceptors.NewAuthHandler(cipherInstance, cfg)
	suite.s = grpc.NewServer(grpc.UnaryInterceptor(interceptorService.UnaryServerInterceptor()))
	pb.RegisterGophkeeperServer(suite.s, suite.server)
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := suite.s.Serve(listen)
		if err != nil {
			log.Fatal(err)
		}
	}()
	suite.client = InitGRPCClient(suite.ctx, &logger, suite.wg, cfg)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestLoginFail() {
	suite.storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("generic_error"))
	code, err := suite.client.Login(modelstorage.RegisterLogin{
		Login:    "some_login",
		Password: "some_password",
	})
	assert.Equal(suite.T(), "rpc error: code = Unauthenticated desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unauthenticated, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestLoginSuccess() {
	suite.storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("some_user_id", nil)
	code, err := suite.client.Login(modelstorage.RegisterLogin{
		Login:    "some_login",
		Password: "some_password",
	})
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestRegisterFail() {
	suite.storage.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	code, err := suite.client.Register(modelstorage.RegisterLogin{
		Login:    "some_login",
		Password: "some_password",
	})
	assert.Equal(suite.T(), "rpc error: code = Unauthenticated desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unauthenticated, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestRegisterSuccess() {
	suite.storage.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	code, err := suite.client.Register(modelstorage.RegisterLogin{
		Login:    "some_login",
		Password: "some_password",
	})
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetTextsBinariesFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	_, code, err := suite.client.GetTextsBinaries()
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetTextsBinariesSuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	storageData := []serverStorage.TextBinaryStorageEntry{
		{
			ID:         0,
			Identifier: suite.cipher.Encode("1"),
			UserID:     suite.cipher.Encode("2"),
			Entry:      suite.cipher.Encode("3"),
			Meta:       suite.cipher.Encode("4"),
		},
	}
	expectedData := make(map[string]modelstorage.TextOrBinary)
	expectedData["1"] = modelstorage.TextOrBinary{
		Identifier: "1",
		Entry:      "3",
		Meta:       "4",
	}
	suite.storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	data, code, err := suite.client.GetTextsBinaries()
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	assert.Equal(suite.T(), expectedData, data)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetLoginsPaswordsFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	_, code, err := suite.client.GetLoginsPasswords()
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetLoginsPaswordsSuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	storageData := []serverStorage.LoginPasswordStorageEntry{
		{
			ID:         0,
			Identifier: suite.cipher.Encode("1"),
			UserID:     suite.cipher.Encode("2"),
			Login:      suite.cipher.Encode("3"),
			Password:   suite.cipher.Encode("4"),
			Meta:       suite.cipher.Encode("5"),
		},
	}
	expectedData := make(map[string]modelstorage.LoginAndPassword)
	expectedData["1"] = modelstorage.LoginAndPassword{
		Identifier: "1",
		Login:      "3",
		Password:   "4",
		Meta:       "5",
	}
	suite.storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	data, code, err := suite.client.GetLoginsPasswords()
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	assert.Equal(suite.T(), expectedData, data)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetBankCardsFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	_, code, err := suite.client.GetBankCards()
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestGetBankCardsSuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	storageData := []serverStorage.BankCardStorageEntry{
		{
			ID:         0,
			Identifier: suite.cipher.Encode("1"),
			UserID:     suite.cipher.Encode("2"),
			Number:     suite.cipher.Encode("3"),
			Holder:     suite.cipher.Encode("4"),
			CVV:        suite.cipher.Encode("5"),
			Meta:       suite.cipher.Encode("6"),
		},
	}
	expectedData := make(map[string]modelstorage.BankCard)
	expectedData["1"] = modelstorage.BankCard{
		Identifier: "1",
		Number:     "3",
		Holder:     "4",
		Cvv:        "5",
		Meta:       "6",
	}
	suite.storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	data, code, err := suite.client.GetBankCards()
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	assert.Equal(suite.T(), expectedData, data)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendBankCardFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetBankCardData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	bankCard := modelstorage.BankCard{
		Identifier: "1",
		Number:     "2",
		Holder:     "3",
		Cvv:        "4",
		Meta:       "5",
	}
	code, err := suite.client.SendBankCard(bankCard)
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendBankCardSuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetBankCardData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	bankCard := modelstorage.BankCard{
		Identifier: "1",
		Number:     "2",
		Holder:     "3",
		Cvv:        "4",
		Meta:       "5",
	}
	code, err := suite.client.SendBankCard(bankCard)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendLoginPasswordFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetLoginPasswordData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	loginPassword := modelstorage.LoginAndPassword{
		Identifier: "1",
		Login:      "2",
		Password:   "3",
		Meta:       "4",
	}
	code, err := suite.client.SendLoginPassword(loginPassword)
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendLoginPasswordSuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetLoginPasswordData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	loginPassword := modelstorage.LoginAndPassword{
		Identifier: "1",
		Login:      "2",
		Password:   "3",
		Meta:       "4",
	}
	code, err := suite.client.SendLoginPassword(loginPassword)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendTextBinaryFail() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetTextBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	textBinary := modelstorage.TextOrBinary{
		Identifier: "1",
		Entry:      "2",
		Meta:       "3",
	}
	code, err := suite.client.SendTextBinary(textBinary)
	assert.Equal(suite.T(), "rpc error: code = Unknown desc = generic_error", err.Error())
	assert.Equal(suite.T(), codes.Unknown, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestSendTextBinarySuccess() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SetTextBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	textBinary := modelstorage.TextOrBinary{
		Identifier: "1",
		Entry:      "2",
		Meta:       "3",
	}
	code, err := suite.client.SendTextBinary(textBinary)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestRemoveBankCard() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	code, err := suite.client.RemoveBankCard("1")
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestRemoveLoginPassword() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	code, err := suite.client.RemoveLoginPassword("1")
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *ClientTestSuite) TestRemoveTextBinary() {
	suite.client.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.client.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.client.token})
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	code, err := suite.client.RemoveTextBinary("1")
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), codes.OK, code)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}
