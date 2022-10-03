package main

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient/client"
	"dk-go-gophkeeper/internal/client/storage/inmemory"
	"dk-go-gophkeeper/internal/client/tui"
	"dk-go-gophkeeper/internal/config"
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
	logger := log.New(flog, `client `, log.LstdFlags|log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.NewDefaultConfiguration()
	err = cfg.Parse()
	if err != nil {
		logger.Fatal(err)
	}
	clientGRPC := client.InitGRPCClient(ctx, logger, wg, cfg)
	storage := inmemory.InitStorage(logger, clientGRPC)
	app := tui.InitTUI(cancel, storage, logger)
	app.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		logger.Print("Shutting down by external signal initiated")
		app.App.Stop()
		cancel()
	}()
	wg.Wait()
}
