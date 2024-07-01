package app

import (
	"github.com/namnv2496/go-wallet/api"
	// "github.com/namnv2496/go-wallet/internal/databaseaccess"
)

type App struct {
	// Database *databaseaccess.Database
	Server *api.Server
}

func NewApp(
	// database *databaseaccess.Database,
	server *api.Server,
) *App {
	return &App{
		// Database: database,
		Server: server,
	}
}
