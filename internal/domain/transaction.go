package domain

type TransactionStatus string

const (
	StatusInitiated          TransactionStatus = "INITIATED"
	StatusPending            TransactionStatus = "PENDING"
	StatusAuthorized         TransactionStatus = "AUTHORIZED"
	StatusProcessing         TransactionStatus = "PROCESSING"
	StatusCompleted          TransactionStatus = "COMPLETED"
	StatusFailed             TransactionStatus = "FAILED"
	StatusReleasedRequested  TransactionStatus = "RELEASED_REQUESTED"
	StatusReleased           TransactionStatus = "RELEASED"
	StatusTimedOut           TransactionStatus = "TIMED_OUT"
	StatusRefunded           TransactionStatus = "REFUNDED"
	StatusBlockRequested     TransactionStatus = "BLOCK_REQUESTED"
	StatusBlocked            TransactionStatus = "BLOCKED"
	StatusNetworkRequested   TransactionStatus = "NETWORK_REQUESTED"
	StatusNetworkConfirmed   TransactionStatus = "NETWORK_CONFIRMED"
	StatusIMPSDebited        TransactionStatus = "IMPS_DEBITED"
	StatusNEFTDebited        TransactionStatus = "NEFT_DEBITED"
	StatusUPIDebited         TransactionStatus = "UPI_DEBITED"
	StatusCardDebited        TransactionStatus = "CARD_DEBITED"
	StatusDebitRequested     TransactionStatus = "DEBIT_REQUESTED"
	StatusIMPSFailed         TransactionStatus = "IMPS_FAILED"
	StatusNEFTFailed         TransactionStatus = "NEFT_FAILED"
	StatusUPIFailed          TransactionStatus = "UPI_FAILED"
	StatusCardFailed         TransactionStatus = "CARD_FAILED"
	StatusNetworkTimedOut    TransactionStatus = "NETWORK_TIMED_OUT"
	StatusNetworkDebitFailed TransactionStatus = "NETWORK_DEBIT_FAILED"
	StatusReleaseeHold       TransactionStatus = "RELEASE_ON_HOLD"
	StatusFinalDebitFromAcc  TransactionStatus = "FINAL_DEBIT_FROM_ACC"
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
type SagaSteps string

const (
	// Initial
	SagaNotStarted SagaSteps = "NOT_STARTED"
	SagaInit       SagaSteps = "INITIATED"

	// Account balance orchestration
	SagaBalanceBlocked SagaSteps = "BALANCE_BLOCK"

	// External network orchestration
	SagaNetworkRequested SagaSteps = "NETWORK_REQUESTED"
	SagaIMPSDebited      SagaSteps = "IMPS_DEBITED"
	SagaNEFTDebited      SagaSteps = "NEFT_DEBITED"
	SagaUPIDebited       SagaSteps = "UPI_DEBITED"
	SagaCardDebited      SagaSteps = "CARD_DEBITED"

	SagaFinalDebitFromAcc SagaSteps = "FINAL_DEBIT_FROM_ACC"

	// Final debit orchestration
	SagaFinalDebited SagaSteps = "FINAL_DEBIT"

	// Compensation flows
	SagaRelease SagaSteps = "RELEASE"

	// Terminal states
	SagaCompleted SagaStatus = "COMPLETED"
	SagaFailed    SagaStatus = "FAILED"
	SagaTimedOut  SagaStatus = "TIMED_OUT"
)
const (
	SagaStatusStarted    string = "STARTED"
	SagaStatusInProgress string = "IN_PROGRESS"
	SagaStatusCompleted  string = "COMPLETED"
	SagaStatusFailed     string = "FAILED"
)

var TransactionStatusToLedger = map[TransactionStatus]LedgerEntryType{

	// Lifecycle
	StatusInitiated: LedgerInit,

	// Authorization / Hold
	StatusAuthorized:     LedgerAuth,
	StatusBlockRequested: LedgerBlock,
	StatusBlocked:        LedgerBlock,

	// Network / Processing
	StatusProcessing:       LedgerProcess,
	StatusNetworkRequested: LedgerProcess,
	StatusNetworkConfirmed: LedgerProcess,

	// Final Debit (by channel)
	StatusIMPSDebited: LedgerDebit,
	StatusNEFTDebited: LedgerDebit,
	StatusUPIDebited:  LedgerDebit,
	StatusCardDebited: LedgerDebit,

	StatusFinalDebitFromAcc: LedgerDebit,

	// Settlement
	StatusCompleted: LedgerSettlement,

	// Compensation
	StatusReleased:     LedgerRelease,
	StatusReleaseeHold: LedgerRelease,

	StatusRefunded: LedgerRefund,

	StatusFailed:   LedgerReversal,
	StatusTimedOut: LedgerReversal,

	StatusIMPSFailed: LedgerReversal,
	StatusNEFTFailed: LedgerReversal,
	StatusUPIFailed:  LedgerReversal,
	StatusCardFailed: LedgerReversal,

	StatusNetworkTimedOut: LedgerReversal,
}

const (
	COMPLETED   string = "COMPLETED"
	IN_PROGRESS string = "IN_PROGRESS"
	FAILED      string = "FAILED"
)
