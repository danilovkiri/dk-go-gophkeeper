package client

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient"
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
)

var (
	_ grpcclient.GRPCClient = (*GRPCClient)(nil)
)

type GRPCClient struct {
	token      string
	md         metadata.MD
	ctx        context.Context
	authorized bool
	conn       *grpc.ClientConn
	logger     *log.Logger
	client     any // here will be pb client interface generated from proto files
}

func InitGRPCClient(ctx context.Context, logger *log.Logger, wg *sync.WaitGroup) *GRPCClient {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	client := GRPCClient{
		token:      "",
		md:         nil,
		ctx:        ctx,
		authorized: false,
		conn:       conn,
		logger:     logger,
		client:     nil,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		err := client.conn.Close()
		if err != nil {
			client.logger.Fatal(err)
		}
	}()
	return &client
}

func (c *GRPCClient) LoginRegister(credentials modelstorage.RegisterLogin) (codes.Code, error) {
	//var header, trailer metadata.MD
	//_, err := c.client.Method(c.ctx, &pb.MethodRequest{}, grpc.Header(&header), grpc.Trailer(&trailer))
	//e, ok := status.FromError(err)
	//if !ok {
	//	return err, e.Code()
	//}
	//token := header.Get("token")
	//md := metadata.New(map[string]string{"token": token[0]})
	//c.token = token[0]
	//c.md = md
	//c.authorized = true
	//return nil, e.Code()
	return 0, nil
}

func (c *GRPCClient) GetTextsBinaries() (map[string]modelstorage.TextOrBinary, codes.Code, error) {
	//newCtx := metadata.NewOutgoingContext(c.ctx, c.md)
	//var request emptypb.Empty
	//resp, err := c.client.Method(newCtx, &request)
	//e, ok := status.FromError(err)
	//if !ok {
	//	return nil, e.Code(), err
	//}
	//// proceed
	return nil, 0, nil
}

func (c *GRPCClient) GetLoginsPasswords() (map[string]modelstorage.LoginAndPassword, codes.Code, error) {
	return nil, 0, nil
}

func (c *GRPCClient) GetBankCards() (map[string]modelstorage.BankCard, codes.Code, error) {
	return nil, 0, nil
}

func (c *GRPCClient) SendBankCard(bankCard modelstorage.BankCard) (codes.Code, error) {
	return 0, nil
}

func (c *GRPCClient) SendLoginPassword(loginPassword modelstorage.LoginAndPassword) (codes.Code, error) {
	return 0, nil
}

func (c *GRPCClient) SendTextBinary(textBinary modelstorage.TextOrBinary) (codes.Code, error) {
	return 0, nil
}

func (c *GRPCClient) Remove(identifier string, db string) (codes.Code, error) {
	return 0, nil
}
