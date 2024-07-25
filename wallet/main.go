package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/app"
	"github.com/namnv2496/go-wallet/internal/wiring"
	"github.com/namnv2496/go-wallet/internal/worker"
)

func main() {

	configConfig, err := config.LoadAllConfig(".")
	if err != nil {
		log.Fatalln("Failed to load config file")
	}
	redisOpt := worker.NewRedisConfigOpt(configConfig)

	app, err := wiring.Initialize(configConfig, redisOpt)

	if err != nil {
		log.Fatalln("Failed to init server: ", err)
	}

	runQueueStart(app)

	// start consumer
	go app.Server.Consumer.Start()

	go app.Server.StartHttpServer(":8080")
	BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func BlockUntilSignal(signals ...os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, signals...)
	<-done
}

func runQueueStart(app *app.App) {
	if err := app.Server.Queue.StartQueue(); err != nil {
		log.Fatalln("Failed to start Distributor task queue")
	}
}
