package app

import (
	"github.com/namnv2496/go-wallet/api/restful"
	// "github.com/namnv2496/go-wallet/internal/databaseaccess"
)

type App struct {
	// Database *databaseaccess.Database
	Server *restful.Server
}

func NewApp(
	// database *databaseaccess.Database,
	server *restful.Server,
) *App {
	return &App{
		// Database: database,
		Server: server,
	}
}
