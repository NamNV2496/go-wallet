package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/namnv2496/go-wallet/internal/wiring"
)

func main() {

	app, err := wiring.Initialize(".")

	if err != nil {
		log.Fatalln("Failed to init server")
	}
	go app.Server.Start(":8080")
	BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func BlockUntilSignal(signals ...os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, signals...)
	<-done
}
