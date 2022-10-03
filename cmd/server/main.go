package main

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/server/api/handlers"
	"dk-go-gophkeeper/internal/server/api/interceptors"
	"dk-go-gophkeeper/internal/server/cipher/v1"
	"dk-go-gophkeeper/internal/server/storage/v1"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildMetadata() {
	// print out build parameters
	switch buildVersion {
	case "":
		fmt.Printf("Build version: %s\n", "N/A")
	default:
		fmt.Printf("Build version: %s\n", buildVersion)
	}
	switch buildDate {
	case "":
		fmt.Printf("Build date: %s\n", "N/A")
	default:
		fmt.Printf("Build date: %s\n", buildDate)
	}
	switch buildCommit {
	case "":
		fmt.Printf("Build commit: %s\n", "N/A")
	default:
		fmt.Printf("Build commit: %s\n", buildCommit)
	}
}

func main() {
	printBuildMetadata()
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()
	logger := log.New(flog, `server `, log.LstdFlags|log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.NewDefaultConfiguration()
	err = cfg.Parse()
	if err != nil {
		logger.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	storageInstance := storage.InitStorage(ctx, logger, cfg, wg)
	server, err := handlers.InitServer(cfg, storageInstance, logger)
	if err != nil {
		logger.Fatal(err)
	}
	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		logger.Fatal(err)
	}
	cipherInstance, err := cipher.NewCipherService(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}
	interceptorService := interceptors.NewAuthHandler(cipherInstance, cfg)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptorService.UnaryServerInterceptor()))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-done
		logger.Print("Server shutdown attempted")
		s.GracefulStop()
		cancel()
	}()
	pb.RegisterGophkeeperServer(s, server)
	logger.Print("Server start attempted")
	if err := s.Serve(listen); err != nil {
		logger.Fatal(err)
	}
	wg.Wait()
	logger.Print("Server shutdown succeeded")
}
