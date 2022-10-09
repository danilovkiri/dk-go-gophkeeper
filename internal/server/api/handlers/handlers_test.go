package handlers

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/mocks"
	"dk-go-gophkeeper/internal/server/api/interceptors"
	"dk-go-gophkeeper/internal/server/cipher/v1"
	serverStorage "dk-go-gophkeeper/internal/server/storage/modelstorage"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"sync"
	"testing"
)

type HandlersTestSuite struct {
	suite.Suite
	storage *mocks.MockDataStorage
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	server  *GophkeeperServer
	s       *grpc.Server
	cfg     *config.Config
	cipher  *cipher.Cipher
	token   string
	md      metadata.MD
}

func (suite *HandlersTestSuite) SetupTest() {
	cfg := config.NewDefaultConfiguration()
	cfg.ServerAddress = ":8080"
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	cfg.HandlersTO = 500
	suite.cfg = cfg
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.wg = &sync.WaitGroup{}
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	suite.storage = mocks.NewMockDataStorage(ctrl)
	server, err := InitServer(cfg, suite.storage, &logger)
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
	suite.wg.Add(1)
	go func() {
		defer suite.wg.Done()
		_ = suite.s.Serve(listen)
	}()
	suite.token = "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	suite.md = metadata.New(map[string]string{suite.cfg.AuthBearerName: suite.token})
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) TestLoginFail1() {
	suite.storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("some_user_id", nil)
	request := pb.LoginRegisterRequest{
		Login:    "some_login",
		Password: "some_password",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.Login(newCtx, &request)
	e, _ := status.FromError(err)
	assert.Equal(suite.T(), codes.Internal, e.Code())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestLoginFail2() {
	suite.storage.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("generic_error"))
	request := pb.LoginRegisterRequest{
		Login:    "some_login",
		Password: "some_password",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.Login(newCtx, &request)
	e, _ := status.FromError(err)
	assert.Equal(suite.T(), codes.Unauthenticated, e.Code())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestRegisterFail1() {
	suite.storage.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	request := pb.LoginRegisterRequest{
		Login:    "some_login",
		Password: "some_password",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.Register(newCtx, &request)
	e, _ := status.FromError(err)
	assert.Equal(suite.T(), codes.Internal, e.Code())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}
func (suite *HandlersTestSuite) TestRegisterFail2() {
	suite.storage.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	request := pb.LoginRegisterRequest{
		Login:    "some_login",
		Password: "some_password",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.Register(newCtx, &request)
	e, _ := status.FromError(err)
	assert.Equal(suite.T(), codes.Unauthenticated, e.Code())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestDeleteBankCard() {
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	request := pb.DeleteBankCardRequest{
		Identifier: "some_id",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.DeleteBankCard(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestDeleteLoginPassword() {
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	request := pb.DeleteLoginPasswordRequest{
		Identifier: "some_id",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.DeleteLoginPassword(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestDeleteTextBinary() {
	suite.storage.EXPECT().SendToQueue(gomock.Any()).Return()
	request := pb.DeleteTextBinaryRequest{
		Identifier: "some_id",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.DeleteTextBinary(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostBankCardSuccess() {
	suite.storage.EXPECT().SetBankCardData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	request := pb.SendBankCardRequest{
		Identifier: "1",
		Number:     "2",
		Holder:     "3",
		Cvv:        "4",
		Meta:       "5",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostBankCard(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostBankCardFail() {
	suite.storage.EXPECT().SetBankCardData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	request := pb.SendBankCardRequest{
		Identifier: "1",
		Number:     "2",
		Holder:     "3",
		Cvv:        "4",
		Meta:       "5",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostBankCard(newCtx, &request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostLoginPasswordSuccess() {
	suite.storage.EXPECT().SetLoginPasswordData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	request := pb.SendLoginPasswordRequest{
		Identifier: "1",
		Login:      "2",
		Password:   "3",
		Meta:       "4",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostLoginPassword(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostLoginPasswordFail() {
	suite.storage.EXPECT().SetLoginPasswordData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	request := pb.SendLoginPasswordRequest{
		Identifier: "1",
		Login:      "2",
		Password:   "3",
		Meta:       "4",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostLoginPassword(newCtx, &request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostTextBinarySuccess() {
	suite.storage.EXPECT().SetTextBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	request := pb.SendTextBinaryRequest{
		Identifier: "1",
		Entry:      "2",
		Meta:       "3",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostTextBinary(newCtx, &request)
	assert.Equal(suite.T(), nil, err)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestPostTextBinaryFail() {
	suite.storage.EXPECT().SetTextBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("generic_error"))
	request := pb.SendTextBinaryRequest{
		Identifier: "1",
		Entry:      "2",
		Meta:       "3",
	}
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	_, err := suite.server.PostTextBinary(newCtx, &request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetBankCardsFail() {
	suite.storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	_, err := suite.server.GetBankCards(newCtx, request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetBankCardsSuccess() {
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
	expResp := pb.GetBankCardsResponse{}
	expSubresp := pb.ResponsePieceBankCard{
		Identifier: "1",
		Number:     "3",
		Holder:     "4",
		Cvv:        "5",
		Meta:       "6",
	}
	expResp.ResponsePiecesBankCards = append(expResp.ResponsePiecesBankCards, &expSubresp)
	suite.storage.EXPECT().GetBankCardData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	resp, err := suite.server.GetBankCards(newCtx, request)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), &expResp, resp)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetLoginsPasswordsFail() {
	suite.storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	_, err := suite.server.GetLoginsPasswords(newCtx, request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetLoginsPasswordsSuccess() {
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
	expResp := pb.GetLoginsPasswordsResponse{}
	expSubresp := pb.ResponsePieceLoginPassword{
		Identifier: "1",
		Login:      "3",
		Password:   "4",
		Meta:       "5",
	}
	expResp.ResponsePiecesLoginsPasswords = append(expResp.GetResponsePiecesLoginsPasswords(), &expSubresp)
	suite.storage.EXPECT().GetLoginPasswordData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	resp, err := suite.server.GetLoginsPasswords(newCtx, request)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), &expResp, resp)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetTextsBinariesFail() {
	suite.storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic_error"))
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	_, err := suite.server.GetTextsBinaries(newCtx, request)
	assert.Equal(suite.T(), "generic_error", err.Error())
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}

func (suite *HandlersTestSuite) TestGetTextsBinariesSuccess() {
	storageData := []serverStorage.TextBinaryStorageEntry{
		{
			ID:         0,
			Identifier: suite.cipher.Encode("1"),
			UserID:     suite.cipher.Encode("2"),
			Entry:      suite.cipher.Encode("3"),
			Meta:       suite.cipher.Encode("4"),
		},
	}
	expResp := pb.GetTextsBinariesResponse{}
	expSubresp := pb.ResponsePieceTextBinary{
		Identifier: "1",
		Entry:      "3",
		Meta:       "4",
	}
	expResp.ResponsePiecesTextsBinaries = append(expResp.ResponsePiecesTextsBinaries, &expSubresp)
	suite.storage.EXPECT().GetTextBinaryData(gomock.Any(), gomock.Any()).Return(storageData, nil)
	newCtx := metadata.NewIncomingContext(context.Background(), suite.md)
	var request *emptypb.Empty
	resp, err := suite.server.GetTextsBinaries(newCtx, request)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), &expResp, resp)
	suite.s.GracefulStop()
	suite.cancel()
	suite.wg.Wait()
}
