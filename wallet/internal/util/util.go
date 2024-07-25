package util

type NewTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"min=1"`
	Currency      string `json:"currency" binding:"required"`
	Message       string `json:"message" binding:"required"`
}
