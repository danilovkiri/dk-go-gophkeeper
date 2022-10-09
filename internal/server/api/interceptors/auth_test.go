package interceptors

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/mocks"
	"dk-go-gophkeeper/internal/server/api/handlers"
	cipher "dk-go-gophkeeper/internal/server/cipher/v1"
	"github.com/rs/zerolog"
	"net"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewAuthHandler(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	_ = NewAuthHandler(cipherInstance, cfg)
}

func TestAuthHandler_AuthFunc_NoMD(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)
	ctx := context.Background()
	err := authHandler.AuthFunc(ctx)
	if e, ok := status.FromError(err); ok {
		assert.Equal(t, codes.Unauthenticated, e.Code())
	} else {
		t.Fatal("Error code was not retrieved")
	}
}

func TestAuthHandler_AuthFunc_CorrectMD(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)
	token := "8773a90a68ebd0fd56dffb1441682414fbec5f454eba9be6129bb00744f50d7f19fd870e97eba101a03b857c675e4836de6f5196"
	md := metadata.New(map[string]string{cfg.AuthBearerName: token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	err := authHandler.AuthFunc(ctx)
	assert.Equal(t, nil, err)
}

func TestAuthHandler_AuthFunc_IncorrectMD(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)
	token := "some_token"
	md := metadata.New(map[string]string{cfg.AuthBearerName: token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	err := authHandler.AuthFunc(ctx)
	if e, ok := status.FromError(err); ok {
		assert.Equal(t, codes.PermissionDenied, e.Code())
	} else {
		t.Fatal("Error code was not retrieved")
	}
}

func TestAuthHandler_AuthFunc_EmptyMD(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)
	md := metadata.New(map[string]string{"some_key": "some_token"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	err := authHandler.AuthFunc(ctx)
	if e, ok := status.FromError(err); ok {
		assert.Equal(t, codes.Unauthenticated, e.Code())
	} else {
		t.Fatal("Error code was not retrieved")
	}
}

func TestAuthHandler_UnaryServerInterceptor_FailDataAccess(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authHandler.UnaryServerInterceptor()))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageInit := mocks.NewMockDataStorage(ctrl)
	server, err := handlers.InitServer(cfg, storageInit, &logger)
	if err != nil {
		t.Fatal(err)
	}
	pb.RegisterGophkeeperServer(s, server)
	go func(t *testing.T) {
		err1 := s.Serve(listen)
		if err1 != nil {
			t.Error(err1)
		}
	}(t)

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// send a request
	ctx := context.Background()
	var header, trailer metadata.MD
	c := pb.NewGophkeeperClient(conn)
	var request emptypb.Empty
	_, err = c.GetBankCards(ctx, &request, grpc.Header(&header), grpc.Trailer(&trailer))
	assert.Equal(t, "rpc error: code = Unauthenticated desc = Empty authorization data was found", err.Error())
	s.GracefulStop()
}

func TestAuthHandler_UnaryServerInterceptor_Login(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authHandler.UnaryServerInterceptor()))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageInit := mocks.NewMockDataStorage(ctrl)
	storageInit.EXPECT().CheckUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("generic_token", nil)
	server, err := handlers.InitServer(cfg, storageInit, &logger)
	if err != nil {
		t.Fatal(err)
	}
	pb.RegisterGophkeeperServer(s, server)
	go func(t *testing.T) {
		err1 := s.Serve(listen)
		if err1 != nil {
			t.Error(err1)
		}
	}(t)

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// send a request
	ctx := context.Background()
	var header, trailer metadata.MD
	c := pb.NewGophkeeperClient(conn)
	var request pb.LoginRegisterRequest
	_, err = c.Login(ctx, &request, grpc.Header(&header), grpc.Trailer(&trailer))
	assert.Equal(t, nil, err)
	s.GracefulStop()
}

func TestAuthHandler_UnaryServerInterceptor_Register(t *testing.T) {
	cfg := config.NewDefaultConfiguration()
	cfg.UserKey = "jds__63h3_7ds"
	cfg.AuthBearerName = "token"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cipherInstance, _ := cipher.NewCipherService(cfg, &logger)
	authHandler := NewAuthHandler(cipherInstance, cfg)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authHandler.UnaryServerInterceptor()))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageInit := mocks.NewMockDataStorage(ctrl)
	storageInit.EXPECT().AddNewUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	server, err := handlers.InitServer(cfg, storageInit, &logger)
	if err != nil {
		t.Fatal(err)
	}
	pb.RegisterGophkeeperServer(s, server)
	go func(t *testing.T) {
		err1 := s.Serve(listen)
		if err1 != nil {
			t.Error(err1)
		}
	}(t)

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// send a request
	ctx := context.Background()
	var header, trailer metadata.MD
	c := pb.NewGophkeeperClient(conn)
	var request pb.LoginRegisterRequest
	_, err = c.Register(ctx, &request, grpc.Header(&header), grpc.Trailer(&trailer))
	assert.Equal(t, nil, err)
	s.GracefulStop()
}
