//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-wallet/api"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/app"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	"github.com/namnv2496/go-wallet/internal/logic"
	"github.com/namnv2496/go-wallet/internal/token"
)

// func Initialize(path string) (*database.Database, error) {
func Initialize(path string) (*app.App, error) {

	wire.Build(
		config.ConfigWireSet,
		databaseaccess.DatabaseWireSet,
		token.TokenWireSet,
		logic.LogicWireSet,
		api.ServerWireSet,
		app.NewApp,
	)
	return nil, nil
}
