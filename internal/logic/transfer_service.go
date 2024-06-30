package logic

import (
	"context"
	"log"

	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

type TransferLogic interface {
	CreateTransfer(ctx context.Context, from int64, to int64, amount int64, currency string, message string) (db.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (db.Transfer, error)
	ListTransfers(ctx context.Context, from int64, to int64, limit int32, offset int32) ([]db.Transfer, error)
	UpdateBalanceOfTransfer(ctx context.Context, from int64, to int64, amount int64) error
	UpdateStatusOfTransfer(ctx context.Context, id int64, status int32, message string) (db.Transfer, error)
}

var _ TransferLogic = (*transerLogic)(nil)

type transerLogic struct {
	database *databaseaccess.Database
}

func NewtranserLogic(
	database *databaseaccess.Database,
) (TransferLogic, error) {
	return &transerLogic{
		database: database,
	}, nil
}

func (t transerLogic) CreateTransfer(
	ctx context.Context,
	from int64,
	to int64,
	amount int64,
	currency string,
	message string,
) (db.Transfer, error) {

	arg := db.CreateTransferParams{
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
		Currency:      currency,
		Status:        0,
		Message:       message,
	}
	return t.database.CreateTransfer(ctx, arg)
}

func (t transerLogic) GetTransfer(
	ctx context.Context,
	id int64,
) (db.Transfer, error) {
	return t.database.GetTransfer(ctx, id)
}

func (t transerLogic) ListTransfers(
	ctx context.Context,
	from int64,
	to int64,
	limit int32,
	offset int32,
) ([]db.Transfer, error) {

	arg := db.ListTransfersParams{
		FromAccountID: from,
		Column2:       to,
		Limit:         limit,
		Offset:        offset,
	}
	return t.database.ListTransfers(ctx, arg)
}

func (t transerLogic) UpdateBalanceOfTransfer(
	ctx context.Context,
	from int64,
	to int64,
	amount int64,
) error {
	// transfer money
	minusArg := db.AddAccountBalanceParams{
		ID:     from,
		Amount: -amount,
	}
	_, err := t.database.AddAccountBalance(ctx, minusArg)
	if err != nil {
		log.Println("Error when update balance: ", err)
		return err
	}

	addArg := db.AddAccountBalanceParams{
		ID:     to,
		Amount: amount,
	}
	_, err = t.database.AddAccountBalance(ctx, addArg)
	if err != nil {
		log.Println("Error when update balance: ", err)
		return err
	}
	return nil
}

func (t transerLogic) UpdateStatusOfTransfer(ctx context.Context, id int64, status int32, message string) (db.Transfer, error) {

	arg := db.UpdateTransferStatusParams{
		ID:      id,
		Status:  status,
		Message: message,
	}
	return t.database.UpdateTransferStatus(ctx, arg)
}
