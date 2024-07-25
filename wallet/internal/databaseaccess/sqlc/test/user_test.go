package db_test

import (
	"context"
	"log"
	"testing"

	"github.com/c2fo/testify/require"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/util"
)

var databaseDatabase *databaseaccess.Database

func init() {
	configConfig, err := config.LoadAllConfig("../../../../")
	if err != nil {
		log.Println("Failed to load config")
	}
	log.Println(configConfig)

	databaseDatabase = databaseaccess.NewDatabase(configConfig)
}

func TestCreateUser(t *testing.T) {

	for i := 0; i < 40; i++ {
		hashedPassword, err := util.HashPassword(util.RandomString(6))
		require.NoError(t, err)

		name := util.RandomOwner()
		arg := db.CreateUserParams{
			Username: pgtype.Text{
				String: name,
				Valid:  true,
			},
			HashedPassword:  hashedPassword,
			FullName:        name,
			Email:           util.RandomEmail(),
			Phone:           util.RandomPhone(10),
			Role:            util.RandomRole(),
			IsEmailVerified: true,
		}

		user, err := databaseDatabase.CreateUser(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, user)
	}
}

func TestGetUser(t *testing.T) {
	arg := db.GetUsersByUsernameOrPhoneParams{
		Username: pgtype.Text{
			String: "%a%",
			Valid:  true,
		},
		Phone: "",
	}
	user, err := databaseDatabase.GetUsersByUsernameOrPhone(
		context.Background(),
		arg,
	)
	if err != nil {
		log.Fatalln("faled to get user")
	}
	log.Println(user)
}
