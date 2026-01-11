package api

type CreateTransactionRequest struct {
	PaymentType string `json:"paymentType"`
	PaymentMode string `json:"paymentMode"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type CreateTransactionResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
