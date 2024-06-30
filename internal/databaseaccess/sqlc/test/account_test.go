package db_test

import (
	"context"
	"log"
	"testing"

	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/util"
)

func init() {
	configConfig, err := config.LoadAllConfig("../../../../")
	if err != nil {
		log.Println("Failed to load config")
	}
	log.Println(configConfig)

	databaseDatabase = databaseaccess.NewDatabase(configConfig)
}

func TestCreateAccount(t *testing.T) {

	for i := 0; i < 40; i++ {
		arg := db.CreateAccountParams{
			UserID:   util.RandomInt(1, 40),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}
		_, err := databaseDatabase.CreateAccount(
			context.Background(),
			arg,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
