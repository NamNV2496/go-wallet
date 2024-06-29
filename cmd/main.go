package main

import (
	"log"

	"github.com/namnv2496/go-wallet/internal/wiring"
)

func main() {

	app, err := wiring.Initialize(".")

	if err != nil {
		log.Fatalln("Failed to init server")
	}
	// consumer.NewConsumer()
	app.Server.Start(":8080")
}
