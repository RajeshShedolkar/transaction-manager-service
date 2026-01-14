package domain

type TransactionStatus string

const (
	StatusInitiated  TransactionStatus = "INITIATED"
	StatusPending    TransactionStatus = "PENDING"
	StatusAuthorized TransactionStatus = "AUTHORIZED"
	StatusCompleted  TransactionStatus = "COMPLETED"
	StatusFailed     TransactionStatus = "FAILED"
	StatusReleased   TransactionStatus = "RELEASED"
	StatusTimedOut   TransactionStatus = "TIMED_OUT"
)

type Transaction struct {
	ID string // TM transaction id

	UserRefId        string
	SourceRefId      string
	DestinationRefId string

	PaymentType string // IMMEDIATE / NEFT / CARD
	PaymentMode string // IMPS / UPI / NEFT / CARD

	Status TransactionStatus

	DcFlag string // D or C

	Amount   int64
	Currency string

	NetworkTxnId string
	GatewayTxnId string
}
