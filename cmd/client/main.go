package main

import (
	"context"
	grpcclient "dk-go-gophkeeper/internal/client/grpcclient/client"
	"dk-go-gophkeeper/internal/client/storage/inmemory"
	"dk-go-gophkeeper/internal/client/tui"
	"dk-go-gophkeeper/internal/config"
	"dk-go-gophkeeper/internal/logger"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	wg := &sync.WaitGroup{}
	flog, err := os.OpenFile(`client.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		loggerInstance.Fatal().Err(err)
	}
	clientGRPC := grpcclient.InitGRPCClient(ctx, loggerInstance, wg, cfg)
	storage := inmemory.InitStorage(loggerInstance, clientGRPC, cfg)
	app := tui.InitTUI(cancel, storage, loggerInstance, cfg)
	app.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		loggerInstance.Warn().Msg("Shutting down by external signal initiated")
		app.App.Stop()
		cancel()
	}()
	wg.Wait()
}
