package restful

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-wallet/internal/mq"
	"github.com/namnv2496/go-wallet/internal/util"
)

func (server *Server) TccTransInTry(ctx *gin.Context) {

	var req util.NewTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("TccTransInTry")
	// check receiver user
	_, err := server.accountService.GetAccount(context.Background(), req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func (server *Server) TccTransInConfirm(ctx *gin.Context) {

	fmt.Println("TccTransInConfirm")
}

func (server *Server) TccTransInCancel(ctx *gin.Context) {

	fmt.Println("TccTransInCancel")
}

func (server *Server) TccTransOutTry(ctx *gin.Context) {

	var req util.NewTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("TccTransOutTry")
	// check authorization
	account, err := server.accountService.GetAccount(context.Background(), req.FromAccountID)
	if err != nil {
		err = fmt.Errorf("cannot find transfer out account_id %d", req.FromAccountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check balance
	if account.Balance < req.Amount {
		err = fmt.Errorf("current money %d is less than request amount %d", account.Balance, req.Amount)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func (server *Server) TccTransOutConfirm(ctx *gin.Context) {

	fmt.Println("TccTransOutConfirm")
	var req util.NewTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.transferService.CreateNewTransfer(
		ctx,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
		req.Currency,
		req.Message,
	)
	if err != nil {
		err = fmt.Errorf("create transfer log error")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Commit
	fmt.Println("Trigger to queue")
	result := &mq.TransferResponse{
		FromId:     req.FromAccountID,
		ToId:       req.ToAccountID,
		TransferId: transfer.ID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Message:    req.Message,
		Status:     1,
	}

	if err := server.producer.SendMessage(
		mq.TopicResult,
		result,
	); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// case send to 3rd
	// if err := server.producer.SendMessage(
	// 	mq.TopicRequest,
	// 	mq.TransferRequest{
	// 		TransferId: transfer.ID,
	// 		FromId:     req.FromAccountID,
	// 		ToId:       req.ToAccountID,
	// 		Amount:     req.Amount,
	// 		Currency:   req.Currency,
	// 		Message:    req.Message,
	// 	},
	// ); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }
	ctx.JSON(http.StatusOK, gin.H{"message": "Transfer successed"})
}

func (server *Server) TccTransOutCancel(ctx *gin.Context) {

	fmt.Println("TccTransOutCancel")
}
