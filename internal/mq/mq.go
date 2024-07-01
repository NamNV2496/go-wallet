package mq

const (
	BrokerList   = "localhost:29092"
	TopicRequest = "go-wallet-transfer-request"
	TopicResult  = "go-wallet-transfer-result"
	ClientId     = "local-client-id"
	Partition    = 0
	MaxRetry     = 1
)

type TransferRequest struct {
	TransferId int64  `json:"transfer_id"`
	FromId     int64  `json:"from_id"`
	ToId       int64  `json:"to_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Message    string `json:"message"`
}

type TransferResponse struct {
	TransferId int64  `json:"transfer_id"`
	FromId     int64  `json:"from_id"`
	ToId       int64  `json:"to_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Message    string `json:"message"`
	Status     int32  `json:"status"`
}
