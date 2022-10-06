// Package grpcclient provides client-side functionality for operating with GRPC server.
package grpcclient

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient"
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
)

// check for interface compliance
var (
	_ grpcclient.GRPCClient = (*GRPCClient)(nil)
)

// GRPCClient defines attributes and methods of a GRPCClient instance.
type GRPCClient struct {
	token  string
	md     metadata.MD
	ctx    context.Context
	conn   *grpc.ClientConn
	logger *log.Logger
	client pb.GophkeeperClient
	cfg    *config.Config
}

// InitGRPCClient initializes GRPCClient instance and listens for context cancellation to close it.
func InitGRPCClient(ctx context.Context, logger *log.Logger, wg *sync.WaitGroup, cfg *config.Config) *GRPCClient {
	logger.Printf("Attempting to initialize GRPC client at %s", cfg.ServerAddress)
	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	client := GRPCClient{
		token:  "",
		md:     nil,
		ctx:    ctx,
		conn:   conn,
		logger: logger,
		client: pb.NewGophkeeperClient(conn),
		cfg:    cfg,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		client.logger.Print("Attempting to close GRPC client")
		err := client.conn.Close()
		if err != nil {
			client.logger.Fatal(err)
		}
		client.logger.Print("GRPC client closed")
	}()
	return &client
}

// Login implements client-side login functionality.
func (c *GRPCClient) Login(credentials modelstorage.RegisterLogin) (codes.Code, error) {
	c.logger.Print("Login attempt received")
	var header, trailer metadata.MD
	_, err := c.client.Login(c.ctx, &pb.LoginRegisterRequest{Login: credentials.Login, Password: credentials.Password}, grpc.Header(&header), grpc.Trailer(&trailer))
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	token := header.Get(c.cfg.AuthBearerName)
	md := metadata.New(map[string]string{c.cfg.AuthBearerName: token[0]})
	c.token = token[0]
	c.md = md
	return e.Code(), nil
}

// Register implements client-side register functionality.
func (c *GRPCClient) Register(credentials modelstorage.RegisterLogin) (codes.Code, error) {
	c.logger.Print("Register attempt received")
	var header, trailer metadata.MD
	_, err := c.client.Register(c.ctx, &pb.LoginRegisterRequest{Login: credentials.Login, Password: credentials.Password}, grpc.Header(&header), grpc.Trailer(&trailer))
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	token := header.Get(c.cfg.AuthBearerName)
	md := metadata.New(map[string]string{c.cfg.AuthBearerName: token[0]})
	c.token = token[0]
	c.md = md
	return e.Code(), nil
}

// GetTextsBinaries implements client-side retrieval of texts/binaries from server and storing them in client storage.
func (c *GRPCClient) GetTextsBinaries() (map[string]modelstorage.TextOrBinary, codes.Code, error) {
	c.logger.Print("Getting texts/binaries attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	var request emptypb.Empty
	resp, err := c.client.GetTextsBinaries(newCtx, &request)
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return nil, e.Code(), err
		}
		return nil, codes.Unknown, err
	}
	result := make(map[string]modelstorage.TextOrBinary)
	for _, responsePiece := range resp.ResponsePiecesTextsBinaries {
		resultPiece := modelstorage.TextOrBinary{
			Identifier: responsePiece.Identifier,
			Entry:      responsePiece.Entry,
			Meta:       responsePiece.Meta,
		}
		result[responsePiece.Identifier] = resultPiece
	}
	return result, e.Code(), nil
}

// GetLoginsPasswords implements client-side retrieval of logins/passwords from server and storing them in client storage.
func (c *GRPCClient) GetLoginsPasswords() (map[string]modelstorage.LoginAndPassword, codes.Code, error) {
	c.logger.Print("Getting logins/passwords attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	var request emptypb.Empty
	resp, err := c.client.GetLoginsPasswords(newCtx, &request)
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return nil, e.Code(), err
		}
		return nil, codes.Unknown, err
	}
	result := make(map[string]modelstorage.LoginAndPassword)
	for _, responsePiece := range resp.ResponsePiecesLoginsPasswords {
		resultPiece := modelstorage.LoginAndPassword{
			Identifier: responsePiece.Identifier,
			Login:      responsePiece.Login,
			Password:   responsePiece.Password,
			Meta:       responsePiece.Meta,
		}
		result[responsePiece.Identifier] = resultPiece
	}
	return result, e.Code(), nil
}

// GetBankCards implements client-side retrieval of bank cards from server and storing them in client storage.
func (c *GRPCClient) GetBankCards() (map[string]modelstorage.BankCard, codes.Code, error) {
	c.logger.Print("Getting bank cards attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	var request emptypb.Empty
	resp, err := c.client.GetBankCards(newCtx, &request)
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return nil, e.Code(), err
		}
		return nil, codes.Unknown, err
	}
	result := make(map[string]modelstorage.BankCard)
	for _, responsePiece := range resp.ResponsePiecesBankCards {
		resultPiece := modelstorage.BankCard{
			Identifier: responsePiece.Identifier,
			Number:     responsePiece.Number,
			Holder:     responsePiece.Holder,
			Cvv:        responsePiece.Cvv,
			Meta:       responsePiece.Meta,
		}
		result[responsePiece.Identifier] = resultPiece
	}
	return result, e.Code(), nil
}

// SendBankCard implements client-side sending of bank card entry to server and client storage.
func (c *GRPCClient) SendBankCard(bankCard modelstorage.BankCard) (codes.Code, error) {
	c.logger.Print("Sending bank card attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.PostBankCard(newCtx, &pb.SendBankCardRequest{Identifier: bankCard.Identifier, Number: bankCard.Number, Holder: bankCard.Holder, Cvv: bankCard.Cvv, Meta: bankCard.Meta})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}

// SendLoginPassword implements client-side sending of login/password entry to server and client storage.
func (c *GRPCClient) SendLoginPassword(loginPassword modelstorage.LoginAndPassword) (codes.Code, error) {
	c.logger.Print("Sending login/password attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.PostLoginPassword(newCtx, &pb.SendLoginPasswordRequest{Identifier: loginPassword.Identifier, Login: loginPassword.Login, Password: loginPassword.Password, Meta: loginPassword.Meta})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}

// SendTextBinary implements client-side sending of text/binary entry to server and client storage.
func (c *GRPCClient) SendTextBinary(textBinary modelstorage.TextOrBinary) (codes.Code, error) {
	c.logger.Print("Sending text/binary attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.PostTextBinary(newCtx, &pb.SendTextBinaryRequest{Identifier: textBinary.Identifier, Entry: textBinary.Entry, Meta: textBinary.Meta})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}

// RemoveBankCard implements client-side removal of bank card entry from server and client storage.
func (c *GRPCClient) RemoveBankCard(identifier string) (codes.Code, error) {
	c.logger.Print("Removing bank card attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.DeleteBankCard(newCtx, &pb.DeleteBankCardRequest{Identifier: identifier})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}

// RemoveLoginPassword implements client-side removal of login/password entry from server and client storage.
func (c *GRPCClient) RemoveLoginPassword(identifier string) (codes.Code, error) {
	c.logger.Print("Removing login/password attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.DeleteLoginPassword(newCtx, &pb.DeleteLoginPasswordRequest{Identifier: identifier})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}

// RemoveTextBinary implements client-side removal of text/binary entry from server and client storage.
func (c *GRPCClient) RemoveTextBinary(identifier string) (codes.Code, error) {
	c.logger.Print("Removing text/binary attempt received")
	newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	_, err := c.client.DeleteTextBinary(newCtx, &pb.DeleteTextBinaryRequest{Identifier: identifier})
	e, ok := status.FromError(err)
	if err != nil {
		c.logger.Print(err)
		if ok {
			return e.Code(), err
		}
		return codes.Unknown, err
	}
	return e.Code(), nil
}
