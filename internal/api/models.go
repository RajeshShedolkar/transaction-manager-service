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

type LedgerEntryResponse struct {
	EntryType string `json:"entryType"`
	Amount    int64  `json:"amount"`
	Source    string `json:"source"`
}

type GetTransactionResponse struct {
	TransactionID string                `json:"transactionId"`
	Status        string                `json:"status"`
	PaymentMode   string                `json:"paymentMode"`
	Amount        int64                 `json:"amount"`
	Currency      string                `json:"currency"`
	Ledger        []LedgerEntryResponse `json:"ledger"`
}
