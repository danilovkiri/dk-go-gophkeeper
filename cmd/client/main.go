package main

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient/client"
	"dk-go-gophkeeper/internal/client/storage/inmemory"
	"dk-go-gophkeeper/internal/client/tui"
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
	clientGRPC := client.InitGRPCClient(ctx, logger, wg)
	storage := inmemory.InitStorage(logger, clientGRPC)
	app := tui.InitTUI(ctx, storage, logger)
	app.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
		app.App.Stop()
	}()

}
