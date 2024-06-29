package logic

import (
	"context"

	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

type TranserLogic interface {
	CreateTransfer(ctx context.Context, from int64, to int64, amount int64, currency string, message string) (db.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (db.Transfer, error)
	ListTransfers(ctx context.Context, from int64, to int64, limit int32, offset int32) ([]db.Transfer, error)
}

var _ TranserLogic = (*transerLogic)(nil)

type transerLogic struct {
	database *databaseaccess.Database
}

func NewtranserLogic(
	database *databaseaccess.Database,
) (TranserLogic, error) {
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
		ToAccountID:   to,
		Limit:         limit,
		Offset:        offset,
	}
	return t.database.ListTransfers(ctx, arg)
}
