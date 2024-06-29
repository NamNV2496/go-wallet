package app

import (
	// "github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/api"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
)

type App struct {
	// Producer sarama.SyncProducer
	Database *databaseaccess.Database
	Server   *api.Server
}

func NewApp(
	// producer sarama.SyncProducer,
	database *databaseaccess.Database,
	server *api.Server,
) *App {
	return &App{
		// Producer: producer,
		Database: database,
		Server:   server,
	}
}
