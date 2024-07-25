package logic

import (
	"context"
	"log"
	"os"

	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/util"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid/v3"
)

var DefaultHTTPServer = "http://localhost:36789/api/dtmsvr"

type TransferLogic interface {
	CreateNewTCCTransfer(ctx context.Context, from int64, to int64, amount int64, currency string, message string) error
	CreateNewTransfer(ctx context.Context, from int64, to int64, amount int64, currency string, message string) (db.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (db.Transfer, error)
	ListTransfers(ctx context.Context, from int64, to int64, limit int32, offset int32) ([]db.Transfer, error)
	UpdateBalanceTransferTx(ctx context.Context, from int64, to int64, amount int64, id int64, status int32, message string) error
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

func (t transerLogic) CreateNewTCCTransfer(
	ctx context.Context,
	from int64,
	to int64,
	amount int64,
	currency string,
	message string,
) error {
	ip := os.Getenv("DefaultHTTPServer")
	if ip == "" {
		ip = DefaultHTTPServer
	}

	gid := shortuuid.New()
	err := dtmcli.TccGlobalTransaction(
		ip,
		gid,
		func(tcc *dtmcli.Tcc) (*resty.Response, error) {
			resp, err := tcc.CallBranch(
				&util.NewTransferRequest{
					FromAccountID: from,
					ToAccountID:   to,
					Amount:        amount,
					Currency:      currency,
					Message:       message,
				},
				"http://localhost:8080/api/v1/TransOut",
				"http://localhost:8080/api/v1/TransOutConfirm",
				"http://localhost:8080/api/v1/TransOutRevert",
			)
			if err != nil {
				return resp, err
			}
			return tcc.CallBranch(
				&util.NewTransferRequest{
					FromAccountID: from,
					ToAccountID:   to,
					Amount:        amount,
					Currency:      currency,
					Message:       message,
				},
				"http://localhost:8080/api/v1/TransIn",
				"http://localhost:8080/api/v1/TransInConfirm",
				"http://localhost:8080/api/v1/TransInRevert",
			)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (t transerLogic) CreateNewTransfer(
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

func (t transerLogic) UpdateBalanceTransferTx(
	ctx context.Context,
	from int64,
	to int64,
	amount int64,
	id int64,
	status int32,
	message string,
) error {

	err := t.database.ExecTx(
		ctx,
		func(query *db.Queries) error {
			// transfer money
			minusArg := db.AddAccountBalanceParams{
				ID:     from,
				Amount: -amount,
			}
			_, err := query.AddAccountBalance(ctx, minusArg)
			if err != nil {
				log.Println("Error when update balance of account: ", from, ", error: ", err)
				return err
			}

			addArg := db.AddAccountBalanceParams{
				ID:     to,
				Amount: amount,
			}
			_, err = query.AddAccountBalance(ctx, addArg)
			if err != nil {
				log.Println("Error when update balance: ", err)
				return err
			}
			arg := db.UpdateTransferStatusParams{
				ID:      id,
				Status:  status,
				Message: message,
			}
			_, err = query.UpdateTransferStatus(ctx, arg)
			if err != nil {
				log.Println("Error when update status of transfer: ", err)
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
