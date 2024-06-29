package logic

import (
	"context"

	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

type AccountLogic interface {
	CreateAccount(ctx context.Context, userID int64, currency string) (db.Account, error)
	GetAccount(ctx context.Context, id int64) (db.Account, error)
	GetAccountsByUserId(ctx context.Context, userID int64) ([]db.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error)
	ListAccounts(ctx context.Context, usetId int64, limit int32, offset int32) ([]db.Account, error)
	UpdateAccount(ctx context.Context, id int64, balance int64) (db.Account, error)
	AddAccountBalance(ctx context.Context, id int64, amount int64) (db.Account, error)
}

var _ AccountLogic = (*accountLogic)(nil)

type accountLogic struct {
	database *databaseaccess.Database
}

func NewAccountLogic(
	database *databaseaccess.Database,
) (AccountLogic, error) {
	return &accountLogic{
		database: database,
	}, nil
}

func (a accountLogic) CreateAccount(ctx context.Context, userID int64, currency string) (db.Account, error) {

	arg := db.CreateAccountParams{
		UserID:   userID,
		Currency: currency,
	}
	return a.database.CreateAccount(ctx, arg)
}

func (a accountLogic) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	return a.database.GetAccount(ctx, id)
}

func (a accountLogic) GetAccountsByUserId(ctx context.Context, userId int64) ([]db.Account, error) {
	return a.database.GetAccountByUserId(ctx, userId)
}

func (a accountLogic) GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error) {
	return db.Account{}, nil
}

func (a accountLogic) ListAccounts(
	ctx context.Context,
	usetId int64,
	limit int32,
	offset int32,
) ([]db.Account, error) {

	arg := db.ListAccountsParams{
		UserID: usetId,
		Limit:  limit,
		Offset: offset,
	}
	return a.database.ListAccounts(ctx, arg)
}

func (a accountLogic) UpdateAccount(ctx context.Context, id int64, balance int64) (db.Account, error) {
	arg := db.UpdateAccountParams{
		ID:      id,
		Balance: balance,
	}
	return a.database.UpdateAccount(ctx, arg)
}

func (a accountLogic) AddAccountBalance(ctx context.Context, id int64, amount int64) (db.Account, error) {
	arg := db.AddAccountBalanceParams{
		ID:     id,
		Amount: amount,
	}
	return a.database.AddAccountBalance(ctx, arg)
}
