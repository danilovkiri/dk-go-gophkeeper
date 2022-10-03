package handlers

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/server/cipher/v1"
	"dk-go-gophkeeper/internal/server/processor"
	service "dk-go-gophkeeper/internal/server/processor/v1"
	"dk-go-gophkeeper/internal/server/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type GophkeeperServer struct {
	pb.UnimplementedGophkeeperServer
	processor processor.Processor
	cfg       *config.Config
	logger    *log.Logger
}

func InitServer(cfg *config.Config, storage storage.DataStorage, logger *log.Logger) (server *GophkeeperServer, err error) {
	logger.Print("Attempting to initialize server")
	cipher, err := cipher.NewCipherService(cfg, logger)
	if err != nil {
		return nil, err
	}
	gophkeeperService := service.InitService(storage, cipher, logger)
	if err != nil {
		return nil, err
	}
	return &GophkeeperServer{processor: gophkeeperService, cfg: cfg, logger: logger}, nil
}

func (s *GophkeeperServer) Register(ctx context.Context, request *pb.LoginRegisterRequest) (*emptypb.Empty, error) {
	s.logger.Print("New register request received")
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	accessToken, err := s.processor.AddNewUser(ctx, request.Login, request.Password)
	if err != nil {
		s.logger.Print("New register request failed")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	md := metadata.New(map[string]string{"token": accessToken})
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		s.logger.Print("New register request failed when sending headers")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	s.logger.Print("New register request succeeded")
	var response emptypb.Empty
	return &response, nil
}

func (s *GophkeeperServer) Login(ctx context.Context, request *pb.LoginRegisterRequest) (*emptypb.Empty, error) {
	s.logger.Print("New login request received")
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	accessToken, err := s.processor.LoginUser(ctx, request.Login, request.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	md := metadata.New(map[string]string{"token": accessToken})
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		s.logger.Print("New login request failed when sending headers")
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.logger.Print("New login request succeeded")
	var response emptypb.Empty
	return &response, nil
}
