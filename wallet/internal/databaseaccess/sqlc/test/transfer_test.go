package db_test

import (
	"context"
	"fmt"
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

func TestCreateTransfer(t *testing.T) {

	for i := 0; i < 40; i++ {
		from := util.RandomInt(1, 40)
		arg := db.CreateTransferParams{
			FromAccountID: from,
			ToAccountID:   util.RandomInt(1, 20),
			Amount:        util.RandomMoney(),
			Status:        0,
			Message:       fmt.Sprintf("%d transfer money", from),
		}
		_, err := databaseDatabase.CreateTransfer(
			context.Background(),
			arg,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
