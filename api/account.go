package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/token"
)

func (server *Server) createAccount(ctx *gin.Context) {

	var req db.CreateAccountParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	account, err := server.accountService.CreateAccount(ctx, authPayload.UserId, req.Currency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, account)
}

type getAccountRequest struct {
	AccountId int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.accountService.GetAccount(context.Background(), req.AccountId)
	if err != nil {
		log.Println("Failed to get account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	if authPayload.UserId != account.UserID {
		err := errors.New("Cannot get account of another person")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccountByUserId(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	account, err := server.accountService.GetAccountsByUserId(context.Background(), authPayload.UserId)
	if err != nil {
		log.Println("Failed to get account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	log.Println("Get accounts of userId: ", authPayload.UserId, req.PageSize, req.PageSize*(req.PageID-1))
	account, err := server.accountService.ListAccounts(ctx, authPayload.UserId, req.PageSize, (req.PageID-1)*req.PageSize)
	if err != nil {
		log.Println("Failed to get account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}
