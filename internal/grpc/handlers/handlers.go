package handlers

//import (
//	"context"
//	"dk-go-gophkeeper/internal/config"
//	pb "dk-go-gophkeeper/internal/grpc/proto"
//	"dk-go-gophkeeper/internal/server/processor"
//	"dk-go-gophkeeper/internal/server/storage"
//)
//
//type GophkeeperServer struct {
//	pb.UnimplementedGophkeeperServer
//	processor processor.Processor
//	cfg       *config.Config
//}
//
//func InitServer(ctx context.Context, cfg *config.Config, storage storage.URLStorage) (server *GophkeeperServer, err error) {
//	gophkeeperService, err := processor.InitShortener(storage)
//	if err != nil {
//		return nil, err
//	}
//	return &GophkeeperServer{processor: shortenerService, cfg: cfg}, nil
//}
