package databaseaccess

import (
	"context"
	"fmt"
	"log"

	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

func (database *Database) ExecTx(
	ctx context.Context,
	fn func(*db.Queries) error,
) error {
	tx, err := database.ConnPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		log.Println("Rollback: ", err)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	log.Println("Commit message")
	return tx.Commit(ctx)
}
