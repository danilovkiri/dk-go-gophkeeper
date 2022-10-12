package main

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	pb "dk-go-gophkeeper/internal/grpc/proto"
	"dk-go-gophkeeper/internal/logger"
	"dk-go-gophkeeper/internal/server/api/handlers"
	"dk-go-gophkeeper/internal/server/api/interceptors"
	cipher "dk-go-gophkeeper/internal/server/cipher/v1"
	storage "dk-go-gophkeeper/internal/server/storage/v1"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

// build parameters to be used with ldflags

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
	loggerInstance := logger.InitLog(flog)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.NewDefaultConfiguration()
	err = cfg.Parse()
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("Config initialization failed")
	}
	wg := &sync.WaitGroup{}
	storageInstance := storage.InitStorage(ctx, loggerInstance, cfg, wg)
	server, err := handlers.InitServer(cfg, storageInstance, loggerInstance)
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("Handlers initialization failed")
	}
	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("Server listening failed")
	}
	cipherInstance, err := cipher.NewCipherService(cfg, loggerInstance)
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("Cipher initialization failed")
	}
	interceptorService := interceptors.NewAuthHandler(cipherInstance, cfg)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptorService.UnaryServerInterceptor()))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-done
		loggerInstance.Warn().Msg("Server shutdown attempted")
		s.GracefulStop()
		cancel()
	}()
	pb.RegisterGophkeeperServer(s, server)
	loggerInstance.Info().Msg("Server start attempted")
	if err := s.Serve(listen); err != nil {
		loggerInstance.Fatal().Err(err).Msg("Server start failed")
	}
	wg.Wait()
	loggerInstance.Info().Msg("Server shutdown succeeded")
}
