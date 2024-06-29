package databaseaccess

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/databaseaccess/migration"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

type Database struct {
	ConnPool *pgxpool.Pool
	*db.Queries
}

var configLocal *config.Config

func NewDatabase(config config.Config) *Database {

	log.Println("Connect to DB: ", config.DBSource)
	ctx := context.Background()
	conPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db: ", err)
	}

	configLocal = &config
	return &Database{
		ConnPool: conPool,
		Queries:  db.New(conPool),
	}
}

func MigrateUp(ctx context.Context) {
	if err := migration.MigrateUp(ctx, *configLocal, false); err != nil {
		log.Fatalln("Failed to run migration!")
	}
}