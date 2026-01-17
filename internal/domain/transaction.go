package domain

type TransactionStatus string

const (
	StatusInitiated         TransactionStatus = "INITIATED"
	StatusPending           TransactionStatus = "PENDING"
	StatusAuthorized        TransactionStatus = "AUTHORIZED"
	StatusProcessing        TransactionStatus = "PROCESSING"
	StatusCompleted         TransactionStatus = "COMPLETED"
	StatusFailed            TransactionStatus = "FAILED"
	StatusReleased          TransactionStatus = "RELEASED"
	StatusTimedOut          TransactionStatus = "TIMED_OUT"
	StatusRefunded          TransactionStatus = "REFUNDED"
	StatusBlockRequested    TransactionStatus = "BLOCK_REQUESTED"
	StatusBlocked           TransactionStatus = "BLOCKED"
	StatusNetworkRequested  TransactionStatus = "NETWORK_REQUESTED"
	StatusNetworkConfirmed  TransactionStatus = "NETWORK_CONFIRMED"
	StatusIMPSDebited       TransactionStatus = "IMPS_DEBITED"
	StatusNEFTDebited       TransactionStatus = "NEFT_DEBITED"
	StatusUPIDebited        TransactionStatus = "UPI_DEBITED"
	StatusCardDebited       TransactionStatus = "CARD_DEBITED"
	StatusIMPSFailed        TransactionStatus = "IMPS_FAILED"
	StatusNEFTFailed        TransactionStatus = "NEFT_FAILED"
	StatusUPIFailed         TransactionStatus = "UPI_FAILED"
	StatusCardFailed        TransactionStatus = "CARD_FAILED"
	StatusNetworkTimedOut   TransactionStatus = "NETWORK_TIMED_OUT"
	StatusReleaseeHold      TransactionStatus = "RELEASE_ON_HOLD"
	StatusFinalDebitFromAcc TransactionStatus = "FINAL_DEBIT_FROM_ACC"
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

type SagaStatus string

const (
	// Initial
	SagaNotStarted SagaStatus = "NOT_STARTED"
	SagaInit       SagaStatus = "INITIATED"

	// Account balance orchestration
	SagaBalanceBlockInProgress SagaStatus = "BALANCE_BLOCK_IN_PROGRESS"
	SagaBalanceBlocked         SagaStatus = "BALANCE_BLOCKED"

	// External network orchestration
	SagaNetworkRequestInProgress SagaStatus = "NETWORK_REQUEST_IN_PROGRESS"
	SagaNetworkConfirmed         SagaStatus = "NETWORK_CONFIRMED"

	// Final debit orchestration
	SagaFinalDebitInProgress SagaStatus = "FINAL_DEBIT_IN_PROGRESS"
	SagaFinalDebited         SagaStatus = "FINAL_DEBIT_COMPLETED"

	// Compensation flows
	SagaReleaseInProgress SagaStatus = "RELEASE_IN_PROGRESS"
	SagaReleased          SagaStatus = "RELEASE_COMPLETED"

	// Terminal states
	SagaCompleted SagaStatus = "COMPLETED"
	SagaFailed    SagaStatus = "FAILED"
	SagaTimedOut  SagaStatus = "TIMED_OUT"
)
const (
	STARTED = "STARTED"
	DONE  = "DONE"
)
