package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-wallet/internal/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"min=1"`
	Currency      string `json:"currency" binding:"required"`
	Message       string `json:"message" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check authorization
	account, err := server.accountService.GetAccount(context.Background(), req.FromAccountID)
	if err != nil {
		log.Println("Failed to get account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	if authPayload.UserId != account.UserID {
		err := errors.New("Cannot transfer money of another person")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check balance
	if account.Balance < req.Amount {
		err := errors.New("Current money is less than request amount")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check receiver user
	_, err = server.accountService.GetAccount(context.Background(), req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// create transfer log
	transfer, err := server.transferService.CreateTransfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount, req.Currency, req.Message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// transfer money
	_, err = server.accountService.AddAccountBalance(ctx, req.FromAccountID, -req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.accountService.AddAccountBalance(ctx, req.ToAccountID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

type getTransferRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(ctx *gin.Context) {

	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	// check authorization
	accounts, err := server.accountService.GetAccountsByUserId(context.Background(), authPayload.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.transferService.GetTransfer(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	for _, account := range accounts {
		if transfer.FromAccountID == account.ID {
			ctx.JSON(http.StatusOK, transfer)
			return
		}
	}
	err = errors.New("Cannot get transfer history of another person")
	ctx.JSON(http.StatusBadRequest, errorResponse(err))
}

type listTransfersRequest struct {
	FromAccountID int64 `form:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `form:"to_account_id"`
	Limit         int32 `form:"limit" binding:"required,min=5"`
	Offset        int32 `form:"offset" binding:"required,min=1"`
}

func (server *Server) listTransfers(ctx *gin.Context) {
	var req listTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check authorization
	account, err := server.accountService.GetAccount(context.Background(), req.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	if authPayload.UserId != account.UserID {
		err := errors.New("Cannot transfer money of another person")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	transfers, err := server.transferService.ListTransfers(ctx, req.FromAccountID, req.ToAccountID, req.Limit, (req.Offset-1)*req.Limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}
