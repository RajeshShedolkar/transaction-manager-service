package api

type CreateTransactionRequest struct {
	UserRefId        string `json:"userRefId"`        // from gateway
	SourceRefId      string `json:"sourceRefId"`      // account/card reference
	DestinationRefId string `json:"destinationRefId"` // merchant/bank reference

	PaymentType string `json:"paymentType"` // IMMEDIATE, NEFT, CARD
	PaymentMode string `json:"paymentMode"` // IMPS, UPI, NEFT, CARD

	DcFlag   string `json:"dcFlag"` // D or C
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`

	NetworkTxnId string `json:"networkTxnId,omitempty"`
	GatewayTxnId string `json:"gatewayTxnId,omitempty"`
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
