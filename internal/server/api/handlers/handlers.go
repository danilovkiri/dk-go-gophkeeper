// Package handlers provides GRPC server-side functionality.
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

// GophkeeperServer defines attributes and methods of a GophkeeperServer instance.
type GophkeeperServer struct {
	pb.UnimplementedGophkeeperServer
	processor processor.Processor
	cfg       *config.Config
	logger    *log.Logger
}

// InitServer initializes a GophkeeperServer instance.
func InitServer(cfg *config.Config, storage storage.DataStorage, logger *log.Logger) (server *GophkeeperServer, err error) {
	logger.Print("Attempting to initialize server")
	cipherInstance, err := cipher.NewCipherService(cfg, logger)
	if err != nil {
		return nil, err
	}
	gophkeeperService := service.InitService(storage, cipherInstance, logger)
	return &GophkeeperServer{processor: gophkeeperService, cfg: cfg, logger: logger}, nil
}

// Register implements server-side register functionality.
func (s *GophkeeperServer) Register(ctx context.Context, request *pb.LoginRegisterRequest) (*emptypb.Empty, error) {
	s.logger.Print("New register request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
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

// Login implements server-side login functionality.
func (s *GophkeeperServer) Login(ctx context.Context, request *pb.LoginRegisterRequest) (*emptypb.Empty, error) {
	s.logger.Print("New login request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
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

// DeleteBankCard performs bank card entry removal from server DB.
func (s *GophkeeperServer) DeleteBankCard(ctx context.Context, request *pb.DeleteBankCardRequest) (*emptypb.Empty, error) {
	s.logger.Print("New DELETE bank card request received")
	userID := s.getUserID(ctx)
	s.processor.Delete(userID, request.Identifier, s.cfg.BankCardDB)
	var response emptypb.Empty
	return &response, nil
}

// DeleteLoginPassword performs login/password entry removal from server DB.
func (s *GophkeeperServer) DeleteLoginPassword(ctx context.Context, request *pb.DeleteLoginPasswordRequest) (*emptypb.Empty, error) {
	s.logger.Print("New DELETE login/password request received")
	userID := s.getUserID(ctx)
	s.processor.Delete(userID, request.Identifier, s.cfg.LoginPasswordDB)
	var response emptypb.Empty
	return &response, nil
}

// DeleteTextBinary performs text/binary entry removal from server DB.
func (s *GophkeeperServer) DeleteTextBinary(ctx context.Context, request *pb.DeleteTextBinaryRequest) (*emptypb.Empty, error) {
	s.logger.Print("New DELETE text/binary request received")
	userID := s.getUserID(ctx)
	s.processor.Delete(userID, request.Identifier, s.cfg.TextBinaryDB)
	var response emptypb.Empty
	return &response, nil
}

// PostBankCard performs bank card entry addition to server DB.
func (s *GophkeeperServer) PostBankCard(ctx context.Context, request *pb.SendBankCardRequest) (*emptypb.Empty, error) {
	s.logger.Print("New POST bank card request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	err := s.processor.SetBankCardData(ctx, userID, request.Identifier, request.Number, request.Holder, request.Cvv, request.Meta)
	if err != nil {
		return nil, err
	}
	var response emptypb.Empty
	return &response, nil
}

// PostLoginPassword performs login/password entry addition to server DB.
func (s *GophkeeperServer) PostLoginPassword(ctx context.Context, request *pb.SendLoginPasswordRequest) (*emptypb.Empty, error) {
	s.logger.Print("New POST login/password request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	err := s.processor.SetLoginPasswordData(ctx, userID, request.Identifier, request.Login, request.Password, request.Meta)
	if err != nil {
		return nil, err
	}
	var response emptypb.Empty
	return &response, nil
}

// PostTextBinary performs text/binary entry addition to server DB.
func (s *GophkeeperServer) PostTextBinary(ctx context.Context, request *pb.SendTextBinaryRequest) (*emptypb.Empty, error) {
	s.logger.Print("New POST text/binary request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	err := s.processor.SetTextBinaryData(ctx, userID, request.Identifier, request.Entry, request.Meta)
	if err != nil {
		return nil, err
	}
	var response emptypb.Empty
	return &response, nil
}

// GetBankCards performs bank card entries retrieval from server DB.
func (s *GophkeeperServer) GetBankCards(ctx context.Context, _ *emptypb.Empty) (*pb.GetBankCardsResponse, error) {
	s.logger.Print("New GET bank cards request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	bankCards, err := s.processor.GetBankCardData(ctx, userID)
	if err != nil {
		return nil, err
	}
	bankCardsResponse := pb.GetBankCardsResponse{}
	for _, piece := range bankCards {
		bankCardResponse := pb.ResponsePieceBankCard{
			Identifier: piece.Identifier,
			Number:     piece.Number,
			Holder:     piece.Holder,
			Cvv:        piece.CVV,
			Meta:       piece.Meta,
		}
		bankCardsResponse.ResponsePiecesBankCards = append(bankCardsResponse.ResponsePiecesBankCards, &bankCardResponse)
	}
	return &bankCardsResponse, nil
}

// GetLoginsPasswords performs login/password entries retrieval from server DB.
func (s *GophkeeperServer) GetLoginsPasswords(ctx context.Context, _ *emptypb.Empty) (*pb.GetLoginsPasswordsResponse, error) {
	s.logger.Print("New GET logins/passwords request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	loginsPasswords, err := s.processor.GetLoginPasswordData(ctx, userID)
	if err != nil {
		return nil, err
	}
	loginsPasswordsResponse := pb.GetLoginsPasswordsResponse{}
	for _, piece := range loginsPasswords {
		loginPasswordResponse := pb.ResponsePieceLoginPassword{
			Identifier: piece.Identifier,
			Login:      piece.Login,
			Password:   piece.Password,
			Meta:       piece.Meta,
		}
		loginsPasswordsResponse.ResponsePiecesLoginsPasswords = append(loginsPasswordsResponse.ResponsePiecesLoginsPasswords, &loginPasswordResponse)
	}
	return &loginsPasswordsResponse, nil
}

// GetTextsBinaries performs text/binary entries retrieval from server DB.
func (s *GophkeeperServer) GetTextsBinaries(ctx context.Context, _ *emptypb.Empty) (*pb.GetTextsBinariesResponse, error) {
	s.logger.Print("New GET texts/binaries request received")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.HandlersTO)*time.Millisecond)
	defer cancel()
	userID := s.getUserID(ctx)
	textsBinaries, err := s.processor.GetTextBinaryData(ctx, userID)
	if err != nil {
		return nil, err
	}
	textsBinariesResponse := pb.GetTextsBinariesResponse{}
	for _, piece := range textsBinaries {
		textBinaryResponse := pb.ResponsePieceTextBinary{
			Identifier: piece.Identifier,
			Entry:      piece.Entry,
			Meta:       piece.Meta,
		}
		textsBinariesResponse.ResponsePiecesTextsBinaries = append(textsBinariesResponse.ResponsePiecesTextsBinaries, &textBinaryResponse)
	}
	return &textsBinariesResponse, nil
}

// getUserID retrieves userID from request context.
func (s *GophkeeperServer) getUserID(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md.Get(s.cfg.AuthBearerName)
	userID := values[0]
	return userID
}
