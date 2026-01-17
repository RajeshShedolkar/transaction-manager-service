package domain

type DebitActKafkaCommand struct {
	Event         string `json:"event"`
	TransactionId string `json:"transactionId"`
	AccountRefId  string `json:"accountRefId"`
	Amount        int64  `json:"amount"`
}
