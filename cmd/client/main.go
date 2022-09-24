package main

import (
	"context"
	"dk-go-gophkeeper/internal/client/storage/inmemory"
	"dk-go-gophkeeper/internal/client/tui"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flog, err := os.OpenFile(`client.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()
	logger := log.New(flog, `client `, log.LstdFlags|log.Lshortfile)
	ctx := context.Background()
	storage := inmemory.InitStorage(logger)
	app := tui.InitTUI(ctx, storage, logger)
	app.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		app.App.Stop()
	}()

}
