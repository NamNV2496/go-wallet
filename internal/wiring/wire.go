//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-wallet/api"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/app"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	"github.com/namnv2496/go-wallet/internal/logic"
	"github.com/namnv2496/go-wallet/internal/mq/consumer"
	"github.com/namnv2496/go-wallet/internal/mq/producer"
	"github.com/namnv2496/go-wallet/internal/token"
)

func Initialize(path string) (*app.App, error) {

	wire.Build(
		config.ConfigWireSet,
		databaseaccess.DatabaseWireSet,
		token.TokenWireSet,
		logic.LogicWireSet,
		producer.NewProducer,
		consumer.NewConsumer,
		api.ServerWireSet,
		app.NewApp,
	)
	return nil, nil
}
