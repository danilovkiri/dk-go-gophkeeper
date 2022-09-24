package main

import (
	"dk-go-gophkeeper/internal/client/tui"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := tui.InitTUI()
	app.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		app.App.Stop()
	}()

}
