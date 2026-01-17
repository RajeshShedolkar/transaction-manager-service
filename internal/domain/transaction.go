package domain

type TransactionStatus string

const (
	StatusInitiated  TransactionStatus = "INITIATED"
	StatusPending    TransactionStatus = "PENDING"
	StatusAuthorized TransactionStatus = "AUTHORIZED"
	StatusProcessing TransactionStatus = "PROCESSING"
	StatusCompleted  TransactionStatus = "COMPLETED"
	StatusFailed     TransactionStatus = "FAILED"
	StatusReleased   TransactionStatus = "RELEASED"
	StatusTimedOut   TransactionStatus = "TIMED_OUT"
	StatusRefunded   TransactionStatus = "REFUNDED"
)

type Transaction struct {
	ID               string // TM transaction id
	UserRefId        string
	SourceRefId      string
	DestinationRefId string
	PaymentType      string // SYNC / ASYNC / IMMEDIATE / DEFFERRED
	PaymentMode      string // IMPS / UPI / NEFT / CARD
	Status           TransactionStatus
	SagaStatus       string
	DcFlag           string // D or C
	Amount           int64
	Currency         string
	NetworkTxnId     string
	GatewayTxnId     string
}

type ConsumerEventForAccountService struct {
	TransactionID string `json:"transactionId"`
	EventID       string `json:"eventId"`
	EventType     string `json:"eventType"`
	UserRefId     string `json:"userRefId"`
	AccountRefId  string `json:"accountRefId"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type ConsumerDLQEventForCardService struct {
	SourceEventBody TxEvent `json:"sourceEventBody"`
	ErrorMessage    string  `json:"errorMessage"`
}
